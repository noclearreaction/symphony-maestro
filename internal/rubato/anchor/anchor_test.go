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
	content := "You are an agent.\n```rubato:anchor\n{\"plugins\":[\"git_status\"]}\n```\nDo your thing."
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block == nil {
		t.Fatal("expected block, got nil")
	}
	if len(block.Plugins) != 1 || block.Plugins[0] != "git_status" {
		t.Errorf("unexpected plugins: %v", block.Plugins)
	}
}

func TestFind_ValidAnchorWithArgs(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[\"git_status\"],\"git_status\":{\"working_dir\":\"/repo\"}}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block == nil {
		t.Fatal("expected block, got nil")
	}
	args := block.Args["git_status"]
	if args == nil {
		t.Fatal("expected git_status args, got nil")
	}
	if args["working_dir"] != "/repo" {
		t.Errorf("unexpected working_dir: %v", args["working_dir"])
	}
}

func TestFind_ValidAnchorMultiplePlugins(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[\"git_status\",\"other\"]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(block.Plugins) != 2 {
		t.Errorf("expected 2 plugins, got %v", block.Plugins)
	}
}

func TestFind_MalformedMissingCloseTag(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[\"git_status\"]}\nno closing fence"
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

func TestFind_MalformedInvalidPluginArgs(t *testing.T) {
	content := "```rubato:anchor\n{\"plugins\":[\"git_status\"],\"git_status\":\"not-an-object\"}\n```"
	_, err := anchor.Find(content)
	if err == nil {
		t.Fatal("expected error for invalid plugin args type")
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
	content := "```rubato:anchor\n{\"plugins\":[\"z_plugin\",\"a_plugin\",\"m_plugin\"]}\n```"
	block, err := anchor.Find(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"z_plugin", "a_plugin", "m_plugin"}
	for i, name := range want {
		if block.Plugins[i] != name {
			t.Errorf("plugin[%d]: want %q got %q", i, name, block.Plugins[i])
		}
	}
}
