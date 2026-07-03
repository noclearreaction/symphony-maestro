package plugin

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/anchor"
)

// gitTimeout is the fixed timeout for all git subprocess calls.
// It is intentionally short: git status on a local repo should complete in
// milliseconds. A 5-second cap prevents slow network filesystems from
// blocking the proxy indefinitely.
const gitTimeout = 5 * time.Second

// GitStatus implements the git_status plugin.
// It reports repository hygiene metrics for the working directory.
type GitStatus struct{}

// NewGitStatus returns a new GitStatus plugin.
func NewGitStatus() *GitStatus { return &GitStatus{} }

func (g *GitStatus) Name() string { return "git_status" }

// Execute runs git commands in the working_dir option (or the process CWD if absent)
// and returns formatted status lines.
func (g *GitStatus) Execute(ctx context.Context, options []anchor.Option) (string, error) {
	dir := ""
	if v, ok := anchor.StringOption(options, "working_dir"); ok {
		dir = v
	}
	ctx, cancel := context.WithTimeout(ctx, gitTimeout)
	defer cancel()
	return gitStatus(ctx, dir)
}

func gitStatus(ctx context.Context, dir string) (string, error) {
	// Confirm this is a git repository (also catches exec failures).
	if _, err := gitOut(ctx, dir, "rev-parse", "--git-dir"); err != nil {
		return "", fmt.Errorf("git_status: not a git repository or git unavailable: %w", err)
	}

	// Bare repositories have no working tree — report as explicit state.
	if out, _ := gitOut(ctx, dir, "rev-parse", "--is-bare-repository"); strings.TrimSpace(out) == "true" {
		return "state: bare", nil
	}

	// Branch or detached-HEAD state.
	var sb strings.Builder
	branch, err := gitOut(ctx, dir, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		// symbolic-ref fails in detached HEAD.
		head, _ := gitOut(ctx, dir, "rev-parse", "--short", "HEAD")
		sb.WriteString("state: detached-head\n")
		sb.WriteString("head: " + strings.TrimSpace(head) + "\n")
	} else {
		sb.WriteString("branch: " + strings.TrimSpace(branch) + "\n")
	}

	// Ahead/behind remote tracking branch (best-effort; zero when no remote).
	ahead, behind := aheadBehind(ctx, dir)
	sb.WriteString(fmt.Sprintf("ahead: %d\n", ahead))
	sb.WriteString(fmt.Sprintf("behind: %d\n", behind))

	// Staged, unstaged tracked-modified, and untracked counts.
	staged, unstaged, untracked := porcelainCounts(ctx, dir)
	sb.WriteString(fmt.Sprintf("staged: %d\n", staged))
	sb.WriteString(fmt.Sprintf("unstaged: %d\n", unstaged))
	sb.WriteString(fmt.Sprintf("untracked: %d", untracked))

	return sb.String(), nil
}

// aheadBehind returns commits ahead and behind the remote tracking branch.
// Returns 0, 0 when no tracking branch is configured.
func aheadBehind(ctx context.Context, dir string) (ahead, behind int) {
	aOut, err := gitOut(ctx, dir, "rev-list", "--count", "@{u}..HEAD")
	if err != nil {
		return 0, 0
	}
	bOut, _ := gitOut(ctx, dir, "rev-list", "--count", "HEAD..@{u}")
	fmt.Sscanf(strings.TrimSpace(aOut), "%d", &ahead)
	fmt.Sscanf(strings.TrimSpace(bOut), "%d", &behind)
	return
}

// porcelainCounts parses `git status --porcelain=v1` output into counts.
//
// XY path lines:
//   - X != ' '/'?' → staged change
//   - Y != ' '/'?' and X != '?' → unstaged tracked change
//   - X == '?' and Y == '?' → untracked file
func porcelainCounts(ctx context.Context, dir string) (staged, unstaged, untracked int) {
	out, err := gitOut(ctx, dir, "status", "--porcelain=v1")
	if err != nil || strings.TrimSpace(out) == "" {
		return
	}
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if len(line) < 2 {
			continue
		}
		x, y := line[0], line[1]
		if x == '?' && y == '?' {
			untracked++
			continue
		}
		if x != ' ' {
			staged++
		}
		if y != ' ' {
			unstaged++
		}
	}
	return
}

// gitOut runs a git command and returns its stdout as a string.
func gitOut(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.Output()
	return string(out), err
}
