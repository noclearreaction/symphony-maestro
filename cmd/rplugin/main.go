package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/anchor"
	"github.com/noclearreaction/symphony-maestro/internal/rubato/plugin"
)

const usageText = `usage: rplugin <plugin-name> [--working-dir <path>] [--args <json>]

Execute a single rubato plugin and print its output to stdout.

Flags:
  --working-dir string   working directory for plugin execution
  --args string          additional plugin args as a JSON object

Available plugins: git_status, go_test`

func main() {
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		fmt.Fprintln(os.Stderr, usageText)
		os.Exit(1)
	}

	pluginName := os.Args[1]

	fs := flag.NewFlagSet("rplugin", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	workingDir := fs.String("working-dir", "", "working directory for plugin execution")
	argsJSON := fs.String("args", "", "additional plugin args as a JSON object")

	if err := fs.Parse(os.Args[2:]); err != nil {
		os.Exit(1)
	}

	// Build args from --args JSON.
	merged := make(map[string]any)
	if *argsJSON != "" {
		if err := json.Unmarshal([]byte(*argsJSON), &merged); err != nil {
			fmt.Fprintf(os.Stderr, "invalid --args JSON: %v\n", err)
			os.Exit(1)
		}
	}

	// --working-dir takes precedence over working_dir in --args.
	if *workingDir != "" {
		merged["working_dir"] = *workingDir
	}

	// Convert args map to []anchor.Option.
	options := make([]anchor.Option, 0, len(merged))
	for name, setting := range merged {
		options = append(options, anchor.Option{Name: name, Setting: setting})
	}

	registry := plugin.NewRegistry(plugin.NewGitStatus(), plugin.NewGoTest())

	results, err := registry.Execute(context.Background(), []anchor.PluginDescriptor{
		{Plugin: pluginName, Options: options},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Print(results[pluginName])
}
