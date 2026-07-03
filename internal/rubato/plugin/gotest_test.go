package plugin_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/anchor"
	"github.com/noclearreaction/symphony-maestro/internal/rubato/plugin"
)

// makeModule creates a temporary Go module directory with the given (filename, content) pairs.
func makeModule(t *testing.T, files ...string) string {
	t.Helper()
	dir := t.TempDir()
	gomod := "module testmod\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i+1 < len(files); i += 2 {
		if err := os.WriteFile(filepath.Join(dir, files[i]), []byte(files[i+1]), 0o644); err != nil {
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
	if !strings.Contains(out, "tests: pass") {
		t.Errorf("expected 'tests: pass' in output:\n%s", out)
	}
	if !strings.Contains(out, "ran:") {
		t.Errorf("expected 'ran:' count in output:\n%s", out)
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
	if !strings.Contains(out, "tests: fail") {
		t.Errorf("expected 'tests: fail' in output:\n%s", out)
	}
	if !strings.Contains(out, "TestAlwaysFail") {
		t.Errorf("expected failure test name in output:\n%s", out)
	}
	if !strings.Contains(out, "intentional failure message") {
		t.Errorf("expected failure message in output:\n%s", out)
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

func TestGoTest_Truncation(t *testing.T) {
	// Write a test that emits 25 log lines then fails — exceeds the 20-line cap.
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
	if !strings.Contains(out, "tests: fail") {
		t.Fatalf("expected 'tests: fail':\n%s", out)
	}
	if !strings.Contains(out, "output truncated") {
		t.Errorf("expected truncation note in output:\n%s", out)
	}
	// The first collected output line is "=== RUN   TestVerboseFail", so the 20
	// slots hold: RUN line + verbose lines 1-19. Line 20 must be truncated.
	if !strings.Contains(out, "verbose output line 19") {
		t.Errorf("expected line 19 to be present (last before truncation), got:\n%s", out)
	}
	if strings.Contains(out, "verbose output line 20") {
		t.Errorf("expected line 20 to be truncated, but found it in output:\n%s", out)
	}
}

func TestGoTest_NonModule(t *testing.T) {
	dir := t.TempDir() // no go.mod
	// A .go file forces the Go toolchain to look for a module — without go.mod it errors.
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := plugin.NewGoTest()
	_, err := p.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err == nil {
		t.Fatal("expected error for non-module directory, got nil")
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
	if !strings.HasPrefix(out, "tests: ") {
		t.Errorf("expected output to start with 'tests: ', got:\n%s", out)
	}
}

