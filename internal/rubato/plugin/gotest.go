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
)

const (
	goTestDefaultTimeout = 60 * time.Second
	goTestMaxFailLines   = 20
)

// GoTest implements the go_test plugin.
// It reports Go unit test results for the working directory.
type GoTest struct{}

// NewGoTest returns a new GoTest plugin.
func NewGoTest() *GoTest { return &GoTest{} }

func (g *GoTest) Name() string { return "go_test" }

// Execute runs go test -json ./... in args["working_dir"] (or process CWD if absent).
// Optional args["timeout_seconds"] (float64) overrides the default 60s timeout.
func (g *GoTest) Execute(ctx context.Context, args map[string]any) (string, error) {
	dir := ""
	if v, ok := args["working_dir"].(string); ok {
		dir = v
	}
	timeout := goTestDefaultTimeout
	if v, ok := args["timeout_seconds"].(float64); ok && v > 0 {
		timeout = time.Duration(v * float64(time.Second))
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return runGoTest(ctx, dir)
}

// testEvent is a single line from go test -json output.
type testEvent struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"`
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
		// No JSON output at all — build failure or not a Go module.
		if len(out) == 0 {
			msg := strings.TrimSpace(stderr.String())
			if msg == "" {
				msg = err.Error()
			}
			return "", fmt.Errorf("go_test: %s", msg)
		}
		// Non-zero exit with output — test failures or build errors; parse below.
	}

	return parseGoTestOutput(out, err)
}

func parseGoTestOutput(raw []byte, execErr error) (string, error) {
	var (
		ran          int
		cachedPkgs   int
		pkgFailed    bool
		failures     []testKey
		testOutputs  = make(map[testKey][]string)
		pkgCached    = make(map[string]bool)
		pkgOutputBuf strings.Builder
	)

	scanner := bufio.NewScanner(bytes.NewReader(raw))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev testEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			continue // skip unparseable lines (e.g., preamble build output)
		}

		key := testKey{ev.Package, ev.Test}

		switch ev.Action {
		case "output":
			if ev.Test != "" {
				testOutputs[key] = append(testOutputs[key], ev.Output)
			} else {
				pkgOutputBuf.WriteString(ev.Output)
				if strings.Contains(ev.Output, "(cached)") {
					pkgCached[ev.Package] = true
				}
			}
		case "pass", "fail", "skip":
			if ev.Test != "" {
				ran++
				if ev.Action == "fail" {
					failures = append(failures, key)
				}
			} else if ev.Action == "fail" {
				pkgFailed = true
			} else if ev.Action == "pass" && pkgCached[ev.Package] {
				cachedPkgs++
			}
		}
	}

	// A package-level failure with no individual test failures means a build
	// or module setup error — surface the package output as an error.
	if pkgFailed && len(failures) == 0 && execErr != nil {
		msg := strings.TrimSpace(pkgOutputBuf.String())
		if msg == "" {
			msg = execErr.Error()
		}
		return "", fmt.Errorf("go_test: %s", msg)
	}

	var sb strings.Builder
	if len(failures) == 0 {
		sb.WriteString("tests: pass\n")
		sb.WriteString(fmt.Sprintf("ran: %d", ran))
		if cachedPkgs > 0 {
			sb.WriteString(fmt.Sprintf("\ncached: %d pkg(s)", cachedPkgs))
		}
	} else {
		sb.WriteString("tests: fail\n")
		sb.WriteString(fmt.Sprintf("ran: %d\n", ran))
		sb.WriteString(fmt.Sprintf("failed: %d", len(failures)))
		for _, f := range failures {
			sb.WriteString("\n--- FAIL: " + f.test)
			lines := testOutputs[f]
			truncated := false
			if len(lines) > goTestMaxFailLines {
				lines = lines[:goTestMaxFailLines]
				truncated = true
			}
			for _, l := range lines {
				sb.WriteString("\n" + strings.TrimRight(l, "\n"))
			}
			if truncated {
				sb.WriteString("\n    ... (output truncated)")
			}
		}
	}
	return sb.String(), nil
}
