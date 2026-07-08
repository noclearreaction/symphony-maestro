package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const upstream = "https://openrouter.ai/api/v1/chat/completions"

var (
	logDir   string
	fileLocks sync.Map // keyed by session key, value *sync.Mutex
)

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Println("warning: OPENROUTER_API_KEY not set — forwarding requests without Authorization header")
	}

	logDir = os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "/logs"
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("failed to create log dir %s: %v", logDir, err)
	}
	log.Printf("writing session logs to %s", logDir)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		handleChatCompletion(w, r, apiKey)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	addr := ":" + port
	log.Printf("openrouter-proxy listening on %s → %s", addr, upstream)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// sessionKey returns the first 8 hex chars of SHA-256(first 512 bytes of messages[0].content).
func sessionKey(messages []map[string]any) string {
	if len(messages) == 0 {
		return "unknown"
	}
	content, _ := messages[0]["content"].(string)
	if len(content) > 512 {
		content = content[:512]
	}
	sum := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", sum[:4]) // 8 hex chars
}

// turnNumber returns the human turn index from the messages array length.
// Turn 1: 2 messages (system+user), turn 2: 4, etc.
func turnNumber(messages []map[string]any) int {
	n := len(messages)
	if n < 2 {
		return 1
	}
	return (n) / 2
}

// fileMutex returns the per-file mutex for a given session key.
func fileMutex(key string) *sync.Mutex {
	mu, _ := fileLocks.LoadOrStore(key, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

// appendLog appends one JSON line to <logDir>/<sessionKey>.ndjson.
func appendLog(key string, entry map[string]any) {
	mu := fileMutex(key)
	mu.Lock()
	defer mu.Unlock()

	line, err := json.Marshal(entry)
	if err != nil {
		log.Printf("log marshal error: %v", err)
		return
	}

	filename := fmt.Sprintf("%s/%s.ndjson", logDir, key)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("log open error: %v", err)
		return
	}
	defer f.Close()
	f.Write(line)
	f.Write([]byte("\n"))
}

func handleChatCompletion(w http.ResponseWriter, r *http.Request, apiKey string) {
	ts := time.Now()

	// Buffer request body so we can log it and forward it
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}

	// Parse messages for session key and turn number
	var reqJSON map[string]any
	var messages []map[string]any
	if json.Unmarshal(reqBody, &reqJSON) == nil {
		if raw, ok := reqJSON["messages"].([]any); ok {
			for _, m := range raw {
				if msg, ok := m.(map[string]any); ok {
					messages = append(messages, msg)
				}
			}
		}
	}
	key := sessionKey(messages)
	turn := turnNumber(messages)

	// Build upstream request
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upstream, bytes.NewReader(reqBody))
	if err != nil {
		http.Error(w, "failed to build upstream request", http.StatusInternalServerError)
		return
	}

	for k, vals := range r.Header {
		switch http.CanonicalHeaderKey(k) {
		case "Authorization", "Content-Length":
		default:
			for _, v := range vals {
				req.Header.Add(k, v)
			}
		}
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		appendLog(key, map[string]any{
			"timestamp": ts.UTC().Format(time.RFC3339Nano),
			"turn":      turn,
			"request":   reqJSON,
			"error":     err.Error(),
		})
		http.Error(w, "upstream request failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("upstream read error: %v", err)
	}

	// Represent response as JSON inline if valid, string otherwise (SSE)
	var respVal any
	if json.Valid(respBody) {
		var respJSON any
		json.Unmarshal(respBody, &respJSON)
		respVal = respJSON
	} else {
		respVal = string(respBody)
	}

	appendLog(key, map[string]any{
		"timestamp": ts.UTC().Format(time.RFC3339Nano),
		"turn":      turn,
		"request":   reqJSON,
		"response":  respVal,
	})

	for k, vals := range resp.Header {
		if http.CanonicalHeaderKey(k) == "Content-Length" {
			continue
		}
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}
