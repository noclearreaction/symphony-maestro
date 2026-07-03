package mutate_test

import (
	"strings"
	"testing"
)

// TestApply_FirstTurn_AllPluginsInjected verifies that on the first turn (no
// prior history) all declared plugins are injected (task 3.6).
func TestApply_FirstTurn_AllPluginsInjected(t *testing.T) {
	outA := "branch: main"
	outB := "tests: pass\nran: 2"
	inj := newTwoPluginInjector("git_status", outA, "go_test", outB)
	sysAnchor := anchorBlock(`{"plugins":[{"plugin":"git_status"},{"plugin":"go_test"}]}`)

	in := body(sysAnchor, "Hello")
	out, err := inj.Apply(t.Context(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])

	if !strings.Contains(lastContent, "```rubato:state") {
		t.Errorf("expected state block on first turn:\n%s", lastContent)
	}
	if !strings.Contains(lastContent, "[git_status]") {
		t.Errorf("expected git_status section:\n%s", lastContent)
	}
	if !strings.Contains(lastContent, "[go_test]") {
		t.Errorf("expected go_test section:\n%s", lastContent)
	}
}

// TestApply_StableTurn_NoStateBlock verifies that when all declared plugins have
// matching prior output, no state block is prepended (task 3.7).
func TestApply_StableTurn_NoStateBlock(t *testing.T) {
	pluginOut := "branch: main\nahead: 0"
	inj := newInjector(pluginOut)
	sysAnchor := anchorBlock(`{"plugins":[{"plugin":"git_status"}]}`)

	stateContent := "```rubato:state\n[git_status]\n" + pluginOut + "\n```\n\nOld question"
	in := bodyWithHistory(sysAnchor, []string{stateContent}, []string{"ack"}, "New question")

	out, err := inj.Apply(t.Context(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])

	if strings.Contains(lastContent, "```rubato:state") {
		t.Errorf("expected no state block on stable turn:\n%s", lastContent)
	}
	if !strings.Contains(lastContent, "New question") {
		t.Errorf("original user content should be preserved:\n%s", lastContent)
	}
}

// TestApply_OnePluginChanges_PartialBlock verifies that when only one plugin
// changes, only that plugin appears in the state block (task 3.8).
func TestApply_OnePluginChanges_PartialBlock(t *testing.T) {
	stableOut := "tests: pass\nran: 3"
	oldGitOut := "branch: main\nahead: 0"
	newGitOut := "branch: feature\nahead: 2"

	// git_status returns newGitOut; go_test returns stableOut.
	inj := newTwoPluginInjector("git_status", newGitOut, "go_test", stableOut)
	sysAnchor := anchorBlock(`{"plugins":[{"plugin":"git_status"},{"plugin":"go_test"}]}`)

	// Prior state: git_status was oldGitOut, go_test was stableOut.
	priorState := "```rubato:state\n[git_status]\n" + oldGitOut + "\n[go_test]\n" + stableOut + "\n```\n\nPrev"
	in := bodyWithHistory(sysAnchor, []string{priorState}, []string{"ack"}, "Next")

	out, err := inj.Apply(t.Context(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])

	if !strings.Contains(lastContent, "```rubato:state") {
		t.Fatalf("expected state block when git_status changed:\n%s", lastContent)
	}
	if !strings.Contains(lastContent, "[git_status]") {
		t.Errorf("expected git_status in state block:\n%s", lastContent)
	}
	// go_test was stable → must NOT appear in the block.
	if strings.Contains(lastContent, "[go_test]") {
		t.Errorf("go_test should not appear when stable:\n%s", lastContent)
	}
}

// TestApply_PluginBeyondMaxAge_Reinjected verifies that a plugin whose last
// injection is beyond the max_age window is re-injected regardless of content
// match (task 3.9).
func TestApply_PluginBeyondMaxAge_Reinjected(t *testing.T) {
	pluginOut := "branch: main"
	inj := newInjector(pluginOut)
	// max_age=1: only look at the immediately preceding message pair.
	sysAnchor := anchorBlock(`{"plugins":[{"plugin":"git_status"}],"options":[{"name":"max_age","setting":1}]}`)

	// State block is 2 turns back (4 prior messages: user, assistant, user, assistant).
	oldState := "```rubato:state\n[git_status]\n" + pluginOut + "\n```\n\nOld"
	in := bodyWithHistory(sysAnchor,
		[]string{oldState, "Middle"},
		[]string{"ack1", "ack2"},
		"Current",
	)

	out, err := inj.Apply(t.Context(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])

	// Even though output is identical, plugin is beyond max_age → re-injected.
	if !strings.Contains(lastContent, "```rubato:state") {
		t.Errorf("expected re-injection when plugin is beyond max_age:\n%s", lastContent)
	}
	if !strings.Contains(lastContent, "[git_status]") {
		t.Errorf("expected git_status section:\n%s", lastContent)
	}
}

// TestApply_MaxAgeZero_AlwaysInjectsAll verifies that max_age=0 injects all
// plugins unconditionally on every turn (task 3.10).
func TestApply_MaxAgeZero_AlwaysInjectsAll(t *testing.T) {
	pluginOut := "branch: main"
	inj := newInjector(pluginOut)
	sysAnchor := anchorBlock(`{"plugins":[{"plugin":"git_status"}],"options":[{"name":"max_age","setting":0}]}`)

	// Even with a matching prior state, max_age=0 forces injection.
	stateContent := "```rubato:state\n[git_status]\n" + pluginOut + "\n```\n\nPrev"
	in := bodyWithHistory(sysAnchor, []string{stateContent}, []string{"ack"}, "Next")

	out, err := inj.Apply(t.Context(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])

	if !strings.Contains(lastContent, "```rubato:state") {
		t.Errorf("expected state block with max_age=0 (always inject):\n%s", lastContent)
	}
}
