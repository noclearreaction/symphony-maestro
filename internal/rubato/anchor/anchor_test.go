package anchor_test

import (
	"testing"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/anchor"
)

func TestFind_NoAnchor(t *testing.T) {
	block, err := anchor.Find("just a regular system message with no anchor")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if block != nil {
		t.Fatalf("expected nil block, got %+v", block)
	}
}

func TestFind_ValidAnchorSinglePlugin(t *testing.T) {
	content := "You are an agent.\n```rubato:anchor\n{\"plugins\":[{\"plugin\":\"git_status\"}]}\n```\nDo your thing."
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block == nil {
		t.Fatal("expected block, got nil")
	}
	if len(block.Plugins) != 1 || block.Plugins[0].Plugin != "git_status" {
		t.Errorf("unexpected plugins: %v", block.Plugins)
	}
	if len(block.Plugins[0].Options) != 0 {
		t.Errorf("expected empty options, got %v", block.Plugins[0].Options)
	}
}

func TestFind_ValidAnchorWithOptions(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[{\"plugin\":\"git_status\",\"options\":[{\"name\":\"working_dir\",\"setting\":\"/repo\"}]}]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block == nil {
		t.Fatal("expected block, got nil")
	}
	if len(block.Plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(block.Plugins))
	}
	opts := block.Plugins[0].Options
	if len(opts) != 1 {
		t.Fatalf("expected 1 option, got %d", len(opts))
	}
	if opts[0].Name != "working_dir" || opts[0].Setting != "/repo" {
		t.Errorf("unexpected option: %v", opts[0])
	}
}

func TestFind_ValidAnchorMultiplePlugins(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[{\"plugin\":\"git_status\"},{\"plugin\":\"go_test\"}]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(block.Plugins) != 2 {
		t.Errorf("expected 2 plugins, got %v", block.Plugins)
	}
}

func TestFind_MalformedMissingCloseTag(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[{\"plugin\":\"git_status\"}]}\nno closing fence"
	_, err := anchor.Find(content)
	if err == nil {
		t.Fatal("expected error for missing close fence")
	}
}

func TestFind_MalformedInvalidJSON(t *testing.T) {
	content := "```rubato:anchor\nnot-json\n```"
	_, err := anchor.Find(content)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFind_MalformedEmptyBody(t *testing.T) {
	content := "```rubato:anchor\n\n```"
	_, err := anchor.Find(content)
	if err == nil {
		t.Fatal("expected error for empty body")
	}
}

func TestFind_MalformedEmptyPlugins(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[]}\n```"
	_, err := anchor.Find(content)
	if err == nil {
		t.Fatal("expected error for empty plugins array")
	}
}

func TestFind_MalformedInvalidPluginDescriptor(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[\"git_status\"]}\n```"
	_, err := anchor.Find(content)
	if err == nil {
		t.Fatal("expected error for string plugin (not a descriptor object)")
	}
}

func TestFind_NoAnchorEmptyString(t *testing.T) {
	block, err := anchor.Find("")
	if err != nil {
		t.Fatalf("expected no error for empty string, got %v", err)
	}
	if block != nil {
		t.Fatalf("expected nil block for empty string")
	}
}

func TestFind_DeclarationOrderPreserved(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[{\"plugin\":\"z_plugin\"},{\"plugin\":\"a_plugin\"},{\"plugin\":\"m_plugin\"}]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"z_plugin", "a_plugin", "m_plugin"}
	for i, name := range want {
		if block.Plugins[i].Plugin != name {
			t.Errorf("plugin[%d]: want %q got %q", i, name, block.Plugins[i].Plugin)
		}
	}
}

// --- new tests for tasks 1.6–1.10 ---

func TestFind_PluginNoOptions_EmptyOptionsList(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[{\"plugin\":\"git_status\"}]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(block.Plugins[0].Options) != 0 {
		t.Errorf("expected empty options, got %v", block.Plugins[0].Options)
	}
}

func TestFind_PluginWithOptions_CorrectValues(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[{\"plugin\":\"go_test\",\"options\":[{\"name\":\"timeout_seconds\",\"setting\":30}]}]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	opts := block.Plugins[0].Options
	if len(opts) != 1 {
		t.Fatalf("expected 1 option, got %d", len(opts))
	}
	// JSON decodes integers as float64 when target is any.
	setting, ok := opts[0].Setting.(float64)
	if !ok {
		t.Fatalf("expected float64 setting, got %T", opts[0].Setting)
	}
	if int(setting) != 30 {
		t.Errorf("expected setting 30, got %v", setting)
	}
}

func TestFind_TopLevelOptionsAbsent_MaxAge100(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[{\"plugin\":\"git_status\"}]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.MaxAge() != 100 {
		t.Errorf("expected MaxAge 100 when absent, got %d", block.MaxAge())
	}
}

func TestFind_TopLevelOptionsMaxAgeZero(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[{\"plugin\":\"git_status\"}],\"options\":[{\"name\":\"max_age\",\"setting\":0}]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.MaxAge() != 0 {
		t.Errorf("expected MaxAge 0, got %d", block.MaxAge())
	}
}

func TestFind_UnknownTopLevelOptionPreserved(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[{\"plugin\":\"git_status\"}],\"options\":[{\"name\":\"unknown_key\",\"setting\":\"val\"}]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(block.Options) != 1 || block.Options[0].Name != "unknown_key" {
		t.Errorf("unknown option not preserved: %v", block.Options)
	}
}

func TestFind_FlagStyleOption_NilSetting(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[{\"plugin\":\"go_test\",\"options\":[{\"name\":\"verbose\"}]}]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	opt := block.Plugins[0].Options[0]
	if opt.Name != "verbose" {
		t.Errorf("expected name 'verbose', got %q", opt.Name)
	}
	if opt.Setting != nil {
		t.Errorf("expected nil setting for flag-style option, got %v", opt.Setting)
	}
}
