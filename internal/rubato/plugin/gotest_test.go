package plugin_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/anchor"
	"github.com/noclearreaction/symphony-maestro/internal/rubato/plugin"
)

// makeModule creates a temporary Go module directory with the given (filename, content) pairs.
// Parent directories are created as needed to support sub-package paths.
func makeModule(t *testing.T, files ...string) string {
	t.Helper()
	dir := t.TempDir()
	gomod := "module testmod\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i+1 < len(files); i += 2 {
		fullPath := filepath.Join(dir, files[i])
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(files[i+1]), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// ---- GoTest plugin tests ----

func TestGoTest_Pass(t *testing.T) {
	dir := makeModule(t, "pass_test.go", `package foo_test
import "testing"
func TestAlwaysPass(t *testing.T) {}
`)
	p := plugin.NewGoTest()
	out, err := p.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "status: pass") {
		t.Errorf("expected 'status: pass' prefix:\n%s", out)
	}
	if !strings.Contains(out, "passed:") {
		t.Errorf("expected 'passed:' in output:\n%s", out)
	}
	if !strings.Contains(out, "failed: 0") {
		t.Errorf("expected 'failed: 0' in output:\n%s", out)
	}
}

func TestGoTest_Fail(t *testing.T) {
	dir := makeModule(t, "fail_test.go", `package foo_test
import "testing"
func TestAlwaysFail(t *testing.T) {
	t.Error("intentional failure message")
}
`)
	p := plugin.NewGoTest()
	out, err := p.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "status: fail") {
		t.Errorf("expected 'status: fail' prefix:\n%s", out)
	}
	if !strings.Contains(out, "test failures:") {
		t.Errorf("expected 'test failures:' section:\n%s", out)
	}
	if !strings.Contains(out, "FAIL TestAlwaysFail") {
		t.Errorf("expected 'FAIL TestAlwaysFail' in output:\n%s", out)
	}
	if !strings.Contains(out, "intentional failure message") {
		t.Errorf("expected failure message in output:\n%s", out)
	}
	if strings.Contains(out, "=== RUN") {
		t.Errorf("framework line '=== RUN' should be stripped:\n%s", out)
	}
}

func TestGoTest_Timeout(t *testing.T) {
	dir := makeModule(t, "slow_test.go", `package foo_test
import (
	"testing"
	"time"
)
func TestSlow(t *testing.T) {
	time.Sleep(30 * time.Second)
}
`)
	p := plugin.NewGoTest()
	_, err := p.Execute(context.Background(), []anchor.Option{
		{Name: "working_dir", Setting: dir},
		{Name: "timeout_seconds", Setting: float64(1)},
	})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timed out") && !strings.Contains(err.Error(), "deadline") {
		t.Errorf("expected timeout-related error message, got: %v", err)
	}
}

func TestGoTest_VerboseFail(t *testing.T) {
	// Write a test that emits 25 log lines then fails — all must appear (no truncation).
	src := `package foo_test
import "testing"
func TestVerboseFail(t *testing.T) {
	for i := 1; i <= 25; i++ {
		t.Logf("verbose output line %d", i)
	}
	t.Error("failure after verbose output")
}
`
	dir := makeModule(t, "verbose_test.go", src)
	p := plugin.NewGoTest()
	out, err := p.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "status: fail") {
		t.Fatalf("expected 'status: fail':\n%s", out)
	}
	// All 25 verbose lines must be present — output is never truncated.
	for i := 1; i <= 25; i++ {
		want := fmt.Sprintf("verbose output line %d", i)
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output (no truncation):\n%s", want, out)
		}
	}
	if !strings.Contains(out, "failure after verbose output") {
		t.Errorf("expected failure message in output:\n%s", out)
	}
	if strings.Contains(out, "=== RUN") {
		t.Errorf("framework line '=== RUN' should be stripped:\n%s", out)
	}
}

func TestGoTest_BuildError(t *testing.T) {
	dir := makeModule(t, "impl_test.go", `package foo_test
import "testing"
func TestUsesUndefined(t *testing.T) {
	_ = UndefinedFunc()
}
`)
	p := plugin.NewGoTest()
	out, err := p.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err != nil {
		t.Fatalf("unexpected hard error (build errors should surface as status: fail): %v", err)
	}
	if !strings.HasPrefix(out, "status: fail") {
		t.Errorf("expected 'status: fail' prefix:\n%s", out)
	}
	if !strings.Contains(out, "build errors:") {
		t.Errorf("expected 'build errors:' section:\n%s", out)
	}
	if !strings.Contains(out, "undefined: UndefinedFunc") {
		t.Errorf("expected 'undefined: UndefinedFunc' in build errors:\n%s", out)
	}
	if !strings.Contains(out, "passed: 0") {
		t.Errorf("expected 'passed: 0':\n%s", out)
	}
}

func TestGoTest_Mixed(t *testing.T) {
	// pkgood has passing tests; pkgbroken has a build error.
	dir := makeModule(t,
		"pkgood/good_test.go", `package pkgood_test
import "testing"
func TestGood(t *testing.T) {}
`,
		"pkgbroken/broken.go", `package pkgbroken
func Oops() { undefinedThing() }
`,
	)
	p := plugin.NewGoTest()
	out, err := p.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err != nil {
		t.Fatalf("unexpected hard error: %v", err)
	}
	if !strings.HasPrefix(out, "status: fail") {
		t.Errorf("expected 'status: fail' prefix:\n%s", out)
	}
	if !strings.Contains(out, "build errors:") {
		t.Errorf("expected 'build errors:' section:\n%s", out)
	}
	if !strings.Contains(out, "undefined: undefinedThing") {
		t.Errorf("expected build error detail:\n%s", out)
	}
	if !strings.Contains(out, "passed: 1") {
		t.Errorf("expected passing test counted (passed: 1):\n%s", out)
	}
}

func TestGoTest_NonModule(t *testing.T) {
	dir := t.TempDir() // no go.mod
	// A .go file forces the Go toolchain to look for a module — without go.mod it errors.
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := plugin.NewGoTest()
	out, err := p.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err != nil {
		t.Fatalf("unexpected hard error (should surface as status: error): %v", err)
	}
	if !strings.HasPrefix(out, "status: error") {
		t.Errorf("expected output to start with 'status: error', got:\n%s", out)
	}
}

func TestGoTest_DefaultCWD(t *testing.T) {
	// Create an isolated module, chdir into it, then run without working_dir.
	// This avoids recursive test execution in the plugin package directory.
	dir := makeModule(t, "pass_test.go", `package foo_test
import "testing"
func TestDefaultPass(t *testing.T) {}
`)
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) }) //nolint:errcheck

	p := plugin.NewGoTest()
	out, err := p.Execute(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error with default CWD: %v", err)
	}
	if !strings.HasPrefix(out, "status: ") {
		t.Errorf("expected output to start with 'status: ', got:\n%s", out)
	}
}
