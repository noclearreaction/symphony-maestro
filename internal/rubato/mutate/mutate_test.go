package mutate_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/mutate"
	"github.com/noclearreaction/symphony-maestro/internal/rubato/plugin"
)

// anchorBlock wraps jsonPayload in the markdown rubato:anchor fence.
func anchorBlock(jsonPayload string) string {
	return "```rubato:anchor\n" + jsonPayload + "\n```"
}

// newInjector returns an Injector wired with a git_status stub plugin.
func newInjector(output string) *mutate.Injector {
	stub := &stubPlugin{name: "git_status", out: output}
	return mutate.NewInjector(plugin.NewRegistry(stub))
}

type stubPlugin struct {
	name string
	out  string
	err  error
}

func (s *stubPlugin) Name() string { return s.name }
func (s *stubPlugin) Execute(_ context.Context, _ map[string]any) (string, error) {
	return s.out, s.err
}

func body(systemContent, userContent string) []byte {
	b, _ := json.Marshal(map[string]any{
		"model": "test-model",
		"messages": []map[string]any{
			{"role": "system", "content": systemContent},
			{"role": "user", "content": userContent},
		},
	})
	return b
}

func bodyExtra(systemContent, userContent string) []byte {
	b, _ := json.Marshal(map[string]any{
		"model":       "test-model",
		"temperature": 0.7,
		"messages": []map[string]any{
			{"role": "system", "content": systemContent},
			{"role": "user", "content": userContent},
		},
	})
	return b
}

func extractMessages(t *testing.T, b []byte) []map[string]json.RawMessage {
	t.Helper()
	var req map[string]json.RawMessage
	if err := json.Unmarshal(b, &req); err != nil {
		t.Fatalf("extractMessages: %v", err)
	}
	var msgs []map[string]json.RawMessage
	if err := json.Unmarshal(req["messages"], &msgs); err != nil {
		t.Fatalf("extractMessages: messages: %v", err)
	}
	return msgs
}

func msgText(t *testing.T, msg map[string]json.RawMessage) string {
	t.Helper()
	var s string
	if err := json.Unmarshal(msg["content"], &s); err != nil {
		t.Fatalf("msgText: %v", err)
	}
	return s
}

func TestApply_NoAnchor_BodyUnchanged(t *testing.T) {
	inj := newInjector("branch: main")
	in := body("You are a helpful assistant.", "Hello!")
	out, err := inj.Apply(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgsIn := extractMessages(t, in)
	msgsOut := extractMessages(t, out)
	if msgText(t, msgsOut[0]) != msgText(t, msgsIn[0]) {
		t.Errorf("messages[0] content changed without anchor")
	}
	if msgText(t, msgsOut[1]) != msgText(t, msgsIn[1]) {
		t.Errorf("messages[-1] content changed without anchor")
	}
}

func TestApply_MalformedAnchor_ReturnsError(t *testing.T) {
	inj := newInjector("")
	in := body("```rubato:anchor\nbad json\n```", "User msg")
	_, err := inj.Apply(context.Background(), in)
	if err == nil {
		t.Fatal("expected error for malformed anchor")
	}
}

func TestApply_ValidAnchor_RuntimeStatePrependedToLastMessage(t *testing.T) {
	pluginOut := "branch: main\nahead: 0\nbehind: 0\nstaged: 0\nunstaged: 0\nuntracked: 0"
	inj := newInjector(pluginOut)
	sysContent := anchorBlock(`{"plugins":["git_status"]}`) + "\nYou are an assistant."
	in := body(sysContent, "What is 2+2?")
	out, err := inj.Apply(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	lastContent := msgText(t, msgs[len(msgs)-1])

	if !strings.Contains(lastContent, "```rubato:state") {
		t.Errorf("runtime-state not prepended to last message:\n%s", lastContent)
	}
	if !strings.Contains(lastContent, "[git_status]") {
		t.Errorf("git_status section missing in state block:\n%s", lastContent)
	}
	if !strings.Contains(lastContent, "What is 2+2?") {
		t.Errorf("original user message lost:\n%s", lastContent)
	}
	if strings.Index(lastContent, "```rubato:state") > strings.Index(lastContent, "What is 2+2?") {
		t.Errorf("runtime-state is after user content, expected before")
	}
}

func TestApply_ValidAnchor_GuidanceInjectedInSystemMessage(t *testing.T) {
	inj := newInjector("branch: main")
	sysContent := anchorBlock(`{"plugins":["git_status"]}`) + " You are an assistant."
	in := body(sysContent, "Hello")
	out, err := inj.Apply(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	sysText := msgText(t, msgs[0])
	if !strings.Contains(sysText, "```rubato:guidance") {
		t.Errorf("guidance block not injected:\n%s", sysText)
	}
	if !strings.Contains(sysText, "git_status") {
		t.Errorf("guidance missing git_status mention:\n%s", sysText)
	}
}

func TestApply_GuidanceIdempotent(t *testing.T) {
	inj := newInjector("branch: main")
	sysContent := anchorBlock(`{"plugins":["git_status"]}`) +
		"\n\n```rubato:guidance\nalready injected\n```"
	in := body(sysContent, "Hello")
	out, err := inj.Apply(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	sysText := msgText(t, msgs[0])
	count := strings.Count(sysText, "```rubato:guidance")
	if count != 1 {
		t.Errorf("expected 1 guidance block, got %d:\n%s", count, sysText)
	}
}

func TestApply_ByteIdenticalGuidance(t *testing.T) {
	inj := newInjector("branch: main")
	sysContent := anchorBlock(`{"plugins":["git_status"]}`) + " System message."

	out1, err := inj.Apply(context.Background(), body(sysContent, "msg1"))
	if err != nil {
		t.Fatalf("run 1: %v", err)
	}
	out2, err := inj.Apply(context.Background(), body(sysContent, "msg2"))
	if err != nil {
		t.Fatalf("run 2: %v", err)
	}

	msgs1 := extractMessages(t, out1)
	msgs2 := extractMessages(t, out2)
	g1 := extractGuidanceBlock(msgText(t, msgs1[0]))
	g2 := extractGuidanceBlock(msgText(t, msgs2[0]))
	if g1 == "" {
		t.Error("guidance block not found in run 1 output")
	}
	if g1 != g2 {
		t.Errorf("guidance not byte-identical:\nrun1=%q\nrun2=%q", g1, g2)
	}
}

func extractGuidanceBlock(s string) string {
	const open = "```rubato:guidance\n"
	const close = "\n```"
	start := strings.Index(s, open)
	if start == -1 {
		return ""
	}
	end := strings.Index(s[start+len(open):], close)
	if end == -1 {
		return ""
	}
	return s[start : start+len(open)+end+len(close)]
}

func TestApply_NoMessages_PassThrough(t *testing.T) {
	inj := newInjector("")
	in := []byte(`{"model":"gpt-4","temperature":0.5}`)
	out, err := inj.Apply(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(in) {
		t.Errorf("body changed without messages field")
	}
}

func TestApply_EmptyMessages_PassThrough(t *testing.T) {
	inj := newInjector("")
	in := []byte(`{"model":"gpt-4","messages":[]}`)
	out, err := inj.Apply(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), `"messages":[]`) {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestApply_SingleMessage_ReturnsError(t *testing.T) {
	inj := newInjector("branch: main")
	in, _ := json.Marshal(map[string]any{
		"model": "test",
		"messages": []map[string]any{
			{"role": "system", "content": anchorBlock(`{"plugins":["git_status"]}`)},
		},
	})
	_, err := inj.Apply(context.Background(), in)
	if err == nil {
		t.Fatal("expected error for single-message request with anchor")
	}
}

func TestApply_UnknownPlugin_ReturnsError(t *testing.T) {
	inj := mutate.NewInjector(plugin.NewRegistry())
	sysContent := anchorBlock(`{"plugins":["no_such_plugin"]}`)
	_, err := inj.Apply(context.Background(), body(sysContent, "Hello"))
	if err == nil {
		t.Fatal("expected error for unknown plugin")
	}
	if !strings.Contains(err.Error(), "no_such_plugin") {
		t.Errorf("error should name the unknown plugin: %v", err)
	}
}

func TestApply_ExtraTopLevelFieldsPreserved(t *testing.T) {
	inj := newInjector("branch: main")
	sysContent := anchorBlock(`{"plugins":["git_status"]}`) + " System."
	in := bodyExtra(sysContent, "User msg")
	out, err := inj.Apply(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var req map[string]json.RawMessage
	if err := json.Unmarshal(out, &req); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	var temp float64
	if err := json.Unmarshal(req["temperature"], &temp); err != nil || temp != 0.7 {
		t.Errorf("temperature field lost or wrong: %v, %v", err, temp)
	}
}

func TestApply_ArrayContent_HandledCorrectly(t *testing.T) {
	inj := newInjector("branch: main")
	in, _ := json.Marshal(map[string]any{
		"model": "test",
		"messages": []map[string]any{
			{"role": "system", "content": anchorBlock(`{"plugins":["git_status"]}`)},
			{"role": "user", "content": []map[string]any{
				{"type": "text", "text": "Array content user message"},
			}},
		},
	})
	out, err := inj.Apply(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := extractMessages(t, out)
	var parts []map[string]json.RawMessage
	if err := json.Unmarshal(msgs[len(msgs)-1]["content"], &parts); err != nil {
		t.Fatalf("expected array content in last message: %v", err)
	}
	if len(parts) < 2 {
		t.Fatalf("expected at least 2 content parts, got %d", len(parts))
	}
	var firstText string
	if err := json.Unmarshal(parts[0]["text"], &firstText); err != nil {
		t.Fatalf("first part text: %v", err)
	}
	if !strings.Contains(firstText, "```rubato:state") {
		t.Errorf("runtime-state not in first content part:\n%s", firstText)
	}
}

func TestApply_InvalidJSON_ReturnsError(t *testing.T) {
	inj := newInjector("")
	_, err := inj.Apply(context.Background(), []byte("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON body")
	}
}
