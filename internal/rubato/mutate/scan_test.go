package mutate_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/mutate"
	"github.com/noclearreaction/symphony-maestro/internal/rubato/plugin"
)

// newTwoPluginInjector returns an Injector with two stub plugins.
func newTwoPluginInjector(nameA, outA, nameB, outB string) *mutate.Injector {
	return mutate.NewInjector(plugin.NewRegistry(
		&stubPlugin{name: nameA, out: outA},
		&stubPlugin{name: nameB, out: outB},
	))
}

// bodyWithHistory builds a multi-turn request body:
//
//	messages[0]  = system (sysContent)
//	For each i:  messages[2i+1] = user (priorUser[i])
//	             messages[2i+2] = assistant (priorAssistant[i])
//	messages[-1] = user (currentUser)
func bodyWithHistory(sysContent string, priorUser, priorAssistant []string, currentUser string) []byte {
	msgs := []map[string]any{
		{"role": "system", "content": sysContent},
	}
	for i := range priorUser {
		msgs = append(msgs, map[string]any{"role": "user", "content": priorUser[i]})
		if i < len(priorAssistant) {
			msgs = append(msgs, map[string]any{"role": "assistant", "content": priorAssistant[i]})
		}
	}
	msgs = append(msgs, map[string]any{"role": "user", "content": currentUser})
	b, _ := json.Marshal(map[string]any{"model": "test-model", "messages": msgs})
	return b
}

// --- Section 2 tests (scanPriorState is unexported; tested via Apply behaviour) ---

// TestScan_SinglePriorBlock verifies that a single prior state block is found
// and used for comparison so that a stable plugin is not re-injected (task 2.3).
func TestScan_SinglePriorBlock(t *testing.T) {
	pluginOut := "branch: main\nahead: 0"
	inj := newInjector(pluginOut)
	sysAnchor := anchorBlock(`{"plugins":[{"plugin":"git_status"}]}`)

	// Prior user message contains a state block with the same output.
	stateContent := "```rubato:state\n[git_status]\n" + pluginOut + "\n```\n\nOld question"
	in := bodyWithHistory(sysAnchor, []string{stateContent}, []string{"ack"}, "Next question")

	out, err := inj.Apply(t.Context(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])

	// Stable turn: scan finds prior output = fresh output → no block prepended.
	if strings.Contains(lastContent, "```rubato:state") {
		t.Errorf("expected no state block when output is stable:\n%s", lastContent)
	}
}

// TestScan_MostRecentWins verifies that when multiple prior state blocks exist,
// the most recent one's output is used for the diff (task 2.4).
func TestScan_MostRecentWins(t *testing.T) {
	newOutput := "branch: feature\nahead: 1"
	oldOutput := "branch: main\nahead: 0"

	// Stub returns newOutput as fresh output.
	inj := newInjector(newOutput)
	sysAnchor := anchorBlock(`{"plugins":[{"plugin":"git_status"}]}`)

	// Two prior turns: oldest has oldOutput, newest has newOutput.
	oldState := "```rubato:state\n[git_status]\n" + oldOutput + "\n```\n\nOld msg"
	newState := "```rubato:state\n[git_status]\n" + newOutput + "\n```\n\nNewer msg"
	in := bodyWithHistory(sysAnchor,
		[]string{oldState, newState},
		[]string{"ack1", "ack2"},
		"Current",
	)

	out, err := inj.Apply(t.Context(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])

	// Most recent prior block has newOutput. Fresh output is also newOutput.
	// → stable → no re-injection.
	if strings.Contains(lastContent, "```rubato:state") {
		t.Errorf("most recent block should win; stable turn should not re-inject:\n%s", lastContent)
	}
}

// TestScan_StopAtMaxAgeBoundary verifies that state blocks beyond the max_age
// window are ignored and the plugin is treated as stale (task 2.5).
func TestScan_StopAtMaxAgeBoundary(t *testing.T) {
	pluginOut := "branch: main"
	inj := newInjector(pluginOut)
	// max_age=1: only look 1 message back from messages[-2].
	sysAnchor := anchorBlock(`{"plugins":[{"plugin":"git_status"}],"options":[{"name":"max_age","setting":1}]}`)

	// State block is 2 turns back (index last-3 relative to last message).
	// With max_age=1, only messages[-2] is scanned — the one directly before current.
	// The state block is in messages[-4], so it's beyond the window.
	oldState := "```rubato:state\n[git_status]\n" + pluginOut + "\n```\n\nOld"
	in := bodyWithHistory(sysAnchor,
		[]string{oldState, "Middle msg"},
		[]string{"ack1", "ack2"},
		"Current",
	)

	out, err := inj.Apply(t.Context(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])

	// State block is beyond max_age=1 window → plugin treated as stale → re-injected.
	if !strings.Contains(lastContent, "```rubato:state") {
		t.Errorf("expected re-injection when prior block is beyond max_age window:\n%s", lastContent)
	}
}

// TestScan_NoPriorBlocks verifies that when no prior state blocks exist,
// an empty scan result causes all plugins to be injected (task 2.6).
func TestScan_NoPriorBlocks(t *testing.T) {
	pluginOut := "branch: main"
	inj := newInjector(pluginOut)
	sysAnchor := anchorBlock(`{"plugins":[{"plugin":"git_status"}]}`)

	// No prior history — first turn.
	in := body(sysAnchor, "Hello")
	out, err := inj.Apply(t.Context(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])

	if !strings.Contains(lastContent, "```rubato:state") {
		t.Errorf("expected state block on first turn (no prior history):\n%s", lastContent)
	}
	if !strings.Contains(lastContent, "[git_status]") {
		t.Errorf("expected git_status section in state block:\n%s", lastContent)
	}
}

// TestScan_MultiPluginAllStable verifies that when two plugins both match their
// prior output, no state block is emitted.
func TestScan_MultiPluginAllStable(t *testing.T) {
	outA := "branch: main"
	outB := "tests: pass\nran: 3"
	inj := newTwoPluginInjector("git_status", outA, "go_test", outB)
	sysAnchor := anchorBlock(`{"plugins":[{"plugin":"git_status"},{"plugin":"go_test"}]}`)

	stateContent := "```rubato:state\n[git_status]\n" + outA + "\n[go_test]\n" + outB + "\n```\n\nPrev"
	in := bodyWithHistory(sysAnchor, []string{stateContent}, []string{"ack"}, "Next")

	out, err := inj.Apply(t.Context(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])
	if strings.Contains(lastContent, "```rubato:state") {
		t.Errorf("expected no state block when all plugins are stable:\n%s", lastContent)
	}
}

