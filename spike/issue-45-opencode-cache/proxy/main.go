package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

const upstream = "https://openrouter.ai/api/v1/chat/completions"

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Println("warning: OPENROUTER_API_KEY not set — forwarding requests without Authorization header")
	}

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

func handleChatCompletion(w http.ResponseWriter, r *http.Request, apiKey string) {
	// Build upstream request with the same body
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upstream, r.Body)
	if err != nil {
		http.Error(w, "failed to build upstream request", http.StatusInternalServerError)
		return
	}

	// Forward relevant request headers; replace Authorization only if key is set
	for key, vals := range r.Header {
		switch http.CanonicalHeaderKey(key) {
		case "Authorization", "Content-Length":
			// Skip: Authorization replaced below (if key set); Content-Length let Go recalculate
		default:
			for _, v := range vals {
				req.Header.Add(key, v)
			}
		}
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "upstream request failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Forward all response headers verbatim, except Content-Length (incompatible with streaming)
	for key, vals := range resp.Header {
		if http.CanonicalHeaderKey(key) == "Content-Length" {
			continue
		}
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Stream response body; flush after each write for SSE
	flusher, canFlush := w.(http.Flusher)
	buf := make([]byte, 4096)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return
			}
			if canFlush {
				flusher.Flush()
			}
		}
		if readErr != nil {
			if readErr != io.EOF {
				log.Printf("upstream read error: %v", readErr)
			}
			return
		}
	}
}
