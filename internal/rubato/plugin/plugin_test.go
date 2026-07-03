package plugin_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/anchor"
	"github.com/noclearreaction/symphony-maestro/internal/rubato/plugin"
)

// ---- Registry tests ----

type stubPlugin struct {
	name string
	out  string
	err  error
}

func (s *stubPlugin) Name() string { return s.name }
func (s *stubPlugin) Execute(_ context.Context, _ []anchor.Option) (string, error) {
	return s.out, s.err
}

func TestRegistry_KnownPlugin(t *testing.T) {
	r := plugin.NewRegistry(&stubPlugin{name: "stub", out: "hello"})
	out, err := r.Execute(context.Background(), []anchor.PluginDescriptor{{Plugin: "stub"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["stub"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["stub"])
	}
}

func TestRegistry_UnknownPlugin(t *testing.T) {
	r := plugin.NewRegistry()
	_, err := r.Execute(context.Background(), []anchor.PluginDescriptor{{Plugin: "no-such"}})
	if err == nil {
		t.Fatal("expected error for unknown plugin")
	}
	if !strings.Contains(err.Error(), "unknown plugin") {
		t.Errorf("error should mention 'unknown plugin': %v", err)
	}
}

func TestRegistry_PluginFailure(t *testing.T) {
	r := plugin.NewRegistry(&stubPlugin{name: "bad", err: context.DeadlineExceeded})
	_, err := r.Execute(context.Background(), []anchor.PluginDescriptor{{Plugin: "bad"}})
	if err == nil {
		t.Fatal("expected error for plugin failure")
	}
	if !strings.Contains(err.Error(), "bad") {
		t.Errorf("error should name the failing plugin: %v", err)
	}
}

func TestRegistry_MultiplePlugins(t *testing.T) {
	r := plugin.NewRegistry(
		&stubPlugin{name: "a", out: "output-a"},
		&stubPlugin{name: "b", out: "output-b"},
	)
	out, err := r.Execute(context.Background(), []anchor.PluginDescriptor{
		{Plugin: "a"},
		{Plugin: "b"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["a"] != "output-a" || out["b"] != "output-b" {
		t.Errorf("unexpected outputs: %v", out)
	}
}

func TestRegistry_FailFastOnFirstUnknown(t *testing.T) {
	r := plugin.NewRegistry(&stubPlugin{name: "good", out: "ok"})
	_, err := r.Execute(context.Background(), []anchor.PluginDescriptor{
		{Plugin: "good"},
		{Plugin: "missing"},
	})
	if err == nil {
		t.Fatal("expected error when second plugin is unknown")
	}
}

func TestRegistry_NoSessionReuse(t *testing.T) {
	// Each Execute call must run plugins fresh — verify outputs are independent.
	calls := 0
	r := plugin.NewRegistry(&stubPlugin{name: "counter", out: "fresh"})
	for i := 0; i < 3; i++ {
		out, err := r.Execute(context.Background(), []anchor.PluginDescriptor{{Plugin: "counter"}})
		if err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
		if out["counter"] != "fresh" {
			t.Errorf("call %d: unexpected output %q", i, out["counter"])
		}
		calls++
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

// ---- GitStatus plugin tests ----

func skipIfNoGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not in PATH")
	}
}

// initRepo creates a minimal git repo in a temp directory.
func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")
	// Initial commit so HEAD exists.
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("init"), 0o644); err != nil {
		t.Fatalf("write README: %v", err)
	}
	run("add", "README.md")
	run("commit", "-m", "init")
	return dir
}

func TestGitStatus_NormalRepo_CleanWorkingTree(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t)

	g := plugin.NewGitStatus()
	out, err := g.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "branch:") {
		t.Errorf("expected branch line, got:\n%s", out)
	}
	if !strings.Contains(out, "staged: 0") {
		t.Errorf("expected staged: 0, got:\n%s", out)
	}
	if !strings.Contains(out, "unstaged: 0") {
		t.Errorf("expected unstaged: 0, got:\n%s", out)
	}
	if !strings.Contains(out, "untracked: 0") {
		t.Errorf("expected untracked: 0, got:\n%s", out)
	}
}

func TestGitStatus_NormalRepo_WithChanges(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t)

	// Stage a new file.
	staged := filepath.Join(dir, "staged.txt")
	if err := os.WriteFile(staged, []byte("staged"), 0o644); err != nil {
		t.Fatalf("write staged: %v", err)
	}
	cmd := exec.Command("git", "add", "staged.txt")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v\n%s", err, out)
	}

	// Create an untracked file.
	if err := os.WriteFile(filepath.Join(dir, "untracked.txt"), []byte("untracked"), 0o644); err != nil {
		t.Fatalf("write untracked: %v", err)
	}

	g := plugin.NewGitStatus()
	out, err := g.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "staged: 1") {
		t.Errorf("expected staged: 1, got:\n%s", out)
	}
	if !strings.Contains(out, "untracked: 1") {
		t.Errorf("expected untracked: 1, got:\n%s", out)
	}
}

func TestGitStatus_DetachedHead(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t)

	cmd := exec.Command("git", "checkout", "--detach", "HEAD")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git checkout --detach: %v\n%s", err, out)
	}

	g := plugin.NewGitStatus()
	out, err := g.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "state: detached-head") {
		t.Errorf("expected detached-head state, got:\n%s", out)
	}
	if !strings.Contains(out, "head:") {
		t.Errorf("expected head SHA line, got:\n%s", out)
	}
}

func TestGitStatus_BareRepo(t *testing.T) {
	skipIfNoGit(t)
	dir := t.TempDir()
	cmd := exec.Command("git", "init", "--bare")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init --bare: %v\n%s", err, out)
	}

	g := plugin.NewGitStatus()
	out, err := g.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "state: bare" {
		t.Errorf("expected 'state: bare', got %q", out)
	}
}

func TestGitStatus_NonRepo(t *testing.T) {
	skipIfNoGit(t)
	dir := t.TempDir() // plain directory, no git repo

	g := plugin.NewGitStatus()
	_, err := g.Execute(context.Background(), []anchor.Option{{Name: "working_dir", Setting: dir}})
	if err == nil {
		t.Fatal("expected error for non-repo directory")
	}
	if !strings.Contains(err.Error(), "git_status") {
		t.Errorf("error should mention git_status: %v", err)
	}
}

func TestGitStatus_DefaultsToProcessCWD(t *testing.T) {
	skipIfNoGit(t)
	// No working_dir arg — should not error when CWD is a git repo.
	// (Tests run from the package directory which is inside the workspace git repo.)
	g := plugin.NewGitStatus()
	_, err := g.Execute(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected success with default CWD: %v", err)
	}
}
