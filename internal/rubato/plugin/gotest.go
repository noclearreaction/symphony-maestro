package plugin

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/anchor"
)

const (
	// goTestDefaultTimeout is the default timeout for go test runs when
	// no timeout_seconds option is provided in the anchor descriptor.
	goTestDefaultTimeout = 60 * time.Second
	// goTestMaxTimeout caps the timeout regardless of the timeout_seconds option
	// to prevent excessively long plugin execution blocking the proxy.
	goTestMaxTimeout = 600 * time.Second // hard cap: 10 minutes
)

// GoTest implements the go_test plugin.
// It reports Go unit test results for the working directory.
type GoTest struct{}

// NewGoTest returns a new GoTest plugin.
func NewGoTest() *GoTest { return &GoTest{} }

func (g *GoTest) Name() string { return "go_test" }

// Execute runs go test -json ./... in the working_dir option (or process CWD if absent).
// Optional timeout_seconds option overrides the default 60s timeout.
func (g *GoTest) Execute(ctx context.Context, options []anchor.Option) (string, error) {
	dir := ""
	if v, ok := anchor.StringOption(options, "working_dir"); ok {
		dir = v
	}
	timeout := goTestDefaultTimeout
	if v := anchor.IntOption(options, "timeout_seconds", 0); v > 0 {
		timeout = time.Duration(v) * time.Second
		if timeout > goTestMaxTimeout {
			timeout = goTestMaxTimeout
		}
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return runGoTest(ctx, dir)
}

// testEvent is a single line from go test -json output.
type testEvent struct {
	Action     string  `json:"Action"`
	ImportPath string  `json:"ImportPath"` // set on build-output/build-fail events
	Package    string  `json:"Package"`
	Test       string  `json:"Test"`
	Output     string  `json:"Output"`
	Elapsed    float64 `json:"Elapsed"`
}

// isFrameworkLine reports whether s is a Go test runner header line that
// carries no diagnostic value (run/pass/fail/skip markers with timing).
func isFrameworkLine(s string) bool {
	return strings.HasPrefix(s, "=== RUN   ") ||
		strings.HasPrefix(s, "=== PAUSE ") ||
		strings.HasPrefix(s, "=== CONT  ") ||
		strings.HasPrefix(s, "--- PASS: ") ||
		strings.HasPrefix(s, "--- FAIL: ") ||
		strings.HasPrefix(s, "--- SKIP: ")
}

// testKey uniquely identifies a test within a package.
type testKey struct{ pkg, test string }

func runGoTest(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "go", "test", "-json", "./...")
	if dir != "" {
		cmd.Dir = dir
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()

	if err != nil {
		// Context cancellation (timeout) takes priority over other errors.
		if ctx.Err() != nil {
			return "", fmt.Errorf("go_test: execution timed out: %w", ctx.Err())
		}
		// No JSON output at all — not a Go module or pre-parse setup failure.
		// Surface as status: error so the AI sees it rather than the proxy
		// failing the whole request.
		if len(out) == 0 {
			msg := strings.TrimSpace(stderr.String())
			if msg == "" {
				msg = err.Error()
			}
			indented := "  " + strings.ReplaceAll(msg, "\n", "\n  ")
			return "status: error\n" + indented + "\n", nil
		}
		// Non-zero exit with output — test failures or build errors; parse below.
	}

	return parseGoTestOutput(out, err)
}

func parseGoTestOutput(raw []byte, execErr error) (string, error) {
	var (
		ran         int
		passed      int
		pkgFailed   bool
		failures    []testKey
		testOutputs = make(map[testKey][]string)
		// Build error tracking: ImportPath -> compiler error lines.
		buildErrors = make(map[string][]string)
		buildPkgs   []string // ordered list of packages with build errors
		buildPkgSet = make(map[string]bool)
		// Package-level output lines for surfacing setup errors.
		setupLines []string
	)

	scanner := bufio.NewScanner(bytes.NewReader(raw))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev testEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			continue // skip unparseable lines
		}

		switch ev.Action {
		case "build-output":
			// Compiler error lines keyed by ImportPath.
			pkg := ev.ImportPath
			if pkg == "" {
				pkg = ev.Package
			}
			content := strings.TrimRight(ev.Output, "\n")
			// Skip "# pkg" header lines — ImportPath already groups them.
			if strings.HasPrefix(content, "#") || content == "" {
				continue
			}
			if !buildPkgSet[pkg] {
				buildPkgSet[pkg] = true
				buildPkgs = append(buildPkgs, pkg)
			}
			buildErrors[pkg] = append(buildErrors[pkg], content)

		case "output":
			if ev.Test != "" {
				// Filter Go test framework header lines — they add noise without
				// diagnostic value (run/pass/fail markers with timing).
				if !isFrameworkLine(ev.Output) {
					key := testKey{ev.Package, ev.Test}
					testOutputs[key] = append(testOutputs[key], ev.Output)
				}
			} else {
				setupLines = append(setupLines, strings.TrimRight(ev.Output, "\n"))
			}

		case "pass":
			if ev.Test != "" {
				ran++
				passed++
			}
		case "fail":
			if ev.Test != "" {
				ran++
				failures = append(failures, testKey{ev.Package, ev.Test})
			} else {
				pkgFailed = true
			}
		// skip: skipped tests do not contribute to ran/passed/failed.
		}
	}

	hasBuildErrors := len(buildPkgs) > 0
	hasTestFailures := len(failures) > 0

	// Setup error: either a package-level failure with no build/test signal at
	// all, or build errors whose import paths are pattern expansions (./...)
	// rather than real packages — both indicate infrastructure failure
	// (module not found, go.mod problems) not a code-level build failure.
	isPatternOnlyBuildError := hasBuildErrors && !hasTestFailures && ran == 0
	for _, pkg := range buildPkgs {
		if !strings.Contains(pkg, "...") {
			isPatternOnlyBuildError = false
			break
		}
	}
	isSetupError := isPatternOnlyBuildError ||
		(pkgFailed && !hasBuildErrors && !hasTestFailures && ran == 0 && execErr != nil)
	if isSetupError {
		var sb strings.Builder
		sb.WriteString("status: error")
		// Collect error lines from build errors (pattern case) or setupLines.
		var errLines []string
		for _, pkg := range buildPkgs {
			errLines = append(errLines, buildErrors[pkg]...)
		}
		for _, l := range setupLines {
			if l == "" || l == "PASS" || l == "FAIL" ||
				strings.HasPrefix(l, "ok") ||
				strings.HasPrefix(l, "FAIL\t") ||
				strings.HasPrefix(l, "#") {
				continue
			}
			errLines = append(errLines, l)
		}
		for _, l := range errLines {
			sb.WriteString("\n  ")
			sb.WriteString(l)
		}
		sb.WriteString("\n")
		return sb.String(), nil
	}

	var sb strings.Builder
	if !hasBuildErrors && !hasTestFailures {
		sb.WriteString("status: pass")
	} else {
		sb.WriteString("status: fail")
	}
	sb.WriteString(fmt.Sprintf("\nran: %d", ran))
	sb.WriteString(fmt.Sprintf("\npassed: %d", passed))
	sb.WriteString(fmt.Sprintf("\nfailed: %d", len(failures)))

	if hasBuildErrors {
		sb.WriteString("\n\nbuild errors:")
		for _, pkg := range buildPkgs {
			sb.WriteString(fmt.Sprintf("\n  %s:", pkg))
			for _, errLine := range buildErrors[pkg] {
				sb.WriteString(fmt.Sprintf("\n    %s", errLine))
			}
		}
	}

	if hasTestFailures {
		sb.WriteString("\n\ntest failures:")
		// Group failures by package preserving encounter order.
		type group struct {
			pkg  string
			keys []testKey
		}
		var groups []group
		pkgIdx := make(map[string]int)
		for _, f := range failures {
			if i, ok := pkgIdx[f.pkg]; ok {
				groups[i].keys = append(groups[i].keys, f)
			} else {
				pkgIdx[f.pkg] = len(groups)
				groups = append(groups, group{pkg: f.pkg, keys: []testKey{f}})
			}
		}
		for _, g := range groups {
			sb.WriteString(fmt.Sprintf("\n  %s:", g.pkg))
			for _, f := range g.keys {
				sb.WriteString(fmt.Sprintf("\n    FAIL %s", f.test))
				for _, l := range testOutputs[f] {
					sb.WriteString("\n" + strings.TrimRight(l, "\n"))
				}
			}
		}
	}

	sb.WriteString("\n")
	return sb.String(), nil
}

