//go:build smoke

package main_test

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// rubatoSubprocessEnv is the sentinel environment variable that causes the
// test binary to run as a rubato server instead of executing the test suite.
const rubatoSubprocessEnv = "RUBATO_TEST_SUBPROCESS"

// smokeRubatoPort is the fixed port used by the smoke test rubato instance.
// Using a fixed port keeps the opencode fixture config static.
// The test fails immediately if the port is already in use.
const smokeRubatoPort = "18080"

func TestSmokeRoundTrip(t *testing.T) {
	// Skip — not fail — on missing external dependencies.
	if os.Getenv("OPENROUTER_API_KEY") == "" {
		t.Skip("OPENROUTER_API_KEY not set")
	}
	if _, err := exec.LookPath("opencode"); err != nil {
		t.Skip("opencode not in PATH")
	}

	rubatoLog := startRubato(t)

	// Unique probe token embedded in the prompt. We check that rubato logged it
	// in the request body — proving this specific request passed through rubato.
	probe := newProbeToken(t)

	fixtureDir, err := filepath.Abs(filepath.Join("testdata", "smoke"))
	if err != nil {
		t.Fatalf("fixture path: %v", err)
	}

	// Run opencode from os.TempDir() so no project opencode.json or .opencode/
	// directory is found — full isolation from repository configuration.
	oc := exec.Command("opencode", "run",
		"--pure",
		"--model", "openrouter/openai/gpt-4o-mini",
		"--agent", "smoke",
		"--format", "json",
		"--title", "smoke",
		"Repeat exactly: "+probe,
	)
	oc.Dir = os.TempDir()
	oc.Env = append(
		envWithout(os.Environ(), "OPENAI_BASE_URL", "OPENAI_API_KEY", "OPENCODE_CONFIG", "OPENCODE_CONFIG_DIR"),
		"OPENCODE_CONFIG="+filepath.Join(fixtureDir, "opencode.json"),
		"OPENCODE_CONFIG_DIR="+fixtureDir,
	)

	out, err := oc.Output()
	if err != nil {
		t.Fatalf("opencode run: %v\noutput:\n%s\nrubato:\n%s", err, out, rubatoLog.String())
	}
	t.Logf("model response: %s", extractResponseText(out))

	// Verify opencode produced at least one valid JSON event — confirms the
	// full round-trip: opencode → rubato → upstream → rubato → opencode.
	if !containsJSONEvent(out) {
		t.Errorf("opencode produced no valid JSON events")
	}

	// Verify the probe token appears in rubato's logged request body — proves
	// this specific request body transited rubato, not just that rubato started.
	if !bytes.Contains(rubatoLog.Bytes(), []byte(probe)) {
		t.Errorf("probe token %q not found in rubato log — request may not have reached rubato\nrubato output:\n%s", probe, rubatoLog.String())
	}
}

// startRubato starts the test binary as a rubato subprocess on smokeRubatoPort,
// waits until it is ready, and registers cleanup via t.Cleanup.
// It fails immediately if the port is already in use.
func startRubato(t *testing.T) *bytes.Buffer {
	t.Helper()

	// Verify the port is free before starting so the error is clear.
	l, err := net.Listen("tcp", "127.0.0.1:"+smokeRubatoPort)
	if err != nil {
		t.Fatalf("smoke test port %s is already in use: %v", smokeRubatoPort, err)
	}
	l.Close()

	t.Logf("rubato port: %s", smokeRubatoPort)

	var rubatoLog bytes.Buffer
	cmd := exec.Command(os.Args[0], "-test.run=^$")
	cmd.Env = append(os.Environ(),
		rubatoSubprocessEnv+"=1",
		"RUBATO_LISTEN_ADDR=:"+smokeRubatoPort,
		"RUBATO_UPSTREAM_URL=https://openrouter.ai/api",
	)
	cmd.Stdout = &rubatoLog
	cmd.Stderr = &rubatoLog

	if err := cmd.Start(); err != nil {
		t.Fatalf("startRubato: %v", err)
	}
	t.Cleanup(func() {
		cmd.Process.Kill()
		cmd.Wait()
	})

	waitReady(t, "127.0.0.1:"+smokeRubatoPort, 5*time.Second)
	t.Logf("rubato ready")

	return &rubatoLog
}

// waitReady polls addr until it accepts a TCP connection or deadline is reached.
func waitReady(t *testing.T, addr string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("service at %s not ready within %s", addr, timeout)
}

// envWithout returns a copy of env with any entries whose key matches one of
// the given keys removed. Use this before appending overrides to prevent the
// parent environment from shadowing them (on Linux, getenv returns the first
// occurrence, so appending alone is not sufficient).
func envWithout(env []string, keys ...string) []string {
	out := make([]string, 0, len(env))
	for _, e := range env {
		key, _, _ := strings.Cut(e, "=")
		skip := false
		for _, k := range keys {
			if key == k {
				skip = true
				break
			}
		}
		if !skip {
			out = append(out, e)
		}
	}
	return out
}

// newProbeToken returns a UUID v4 string prefixed with "rubato-probe-".
func newProbeToken(t *testing.T) string {
	t.Helper()
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatalf("newProbeToken: %v", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("rubato-probe-%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// extractResponseText concatenates the text content from all "text" type
// JSON events in opencode's output — the model's actual reply.
func extractResponseText(output []byte) string {
	var sb strings.Builder
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		var ev struct {
			Type string `json:"type"`
			Part struct {
				Text string `json:"text"`
			} `json:"part"`
		}
		if json.Unmarshal(line, &ev) == nil && ev.Type == "text" {
			sb.WriteString(ev.Part.Text)
		}
	}
	return sb.String()
}

// containsJSONEvent reports whether output contains at least one line that
// deserialises as a JSON object.
func containsJSONEvent(output []byte) bool {
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		var v map[string]any
		if json.Unmarshal(line, &v) == nil {
			return true
		}
	}
	return false
}
