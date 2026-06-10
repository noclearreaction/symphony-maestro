## 1. Environment Setup

- [x] 1.1 Install Go 1.26.4 on Ubuntu/WSL: from https://go.dev/dl/ (verify: `go version` shows go1.26.4)
- [x] 1.2 Verify Docker and `docker network` commands are available: `docker --version && docker network ls`
- [x] 1.3 Create Docker user-defined network for spike containers: `docker network create spike-net`

## 2. Proxy Scaffold

- [x] 2.1 Create `spike/issue-45-opencode-cache/proxy/` directory
- [x] 2.2 Initialize Go module: `go mod init github.com/noclearreaction/symphony-director/spike/proxy` inside `proxy/`
- [x] 2.3 Create `proxy/main.go` with HTTP server skeleton: `net/http` listener on `PORT` env var (default `8080`), single route `POST /v1/chat/completions`, all other routes return 404

## 3. Core Proxy Logic

- [x] 3.1 Implement startup behavior: log a warning if `OPENROUTER_API_KEY` is not set (key is optional — OpenRouter accepts keyless requests for free models); set `Authorization` header only when key is present
- [x] 3.2 Implement request forwarding: copy incoming request body, set `Authorization: Bearer ${OPENROUTER_API_KEY}`, forward to `https://openrouter.ai/api/v1/chat/completions`
- [x] 3.3 Implement SSE streaming passthrough: forward all response headers verbatim (strip `Content-Length`), copy response body using `io.Copy` with `http.Flusher.Flush()` after each write
- [x] 3.4 Verify non-streaming response path also works (full JSON body forwarded correctly — confirmed via curl test)

## 4. Docker Container

- [x] 4.1 Create `proxy/Dockerfile`: multi-stage build — `golang:1.26` builder stage compiles static binary (`CGO_ENABLED=0`), `alpine` runtime stage copies binary
- [x] 4.2 Build proxy image: `docker build -t openrouter-proxy spike/issue-45-opencode-cache/proxy/`
- [x] 4.3 Verify image starts correctly without `OPENROUTER_API_KEY` (logs warning, continues — does not exit)

## 5. Fixture Integration

- [x] 5.1 Add custom provider entry to `spike/issue-45-opencode-cache/fixture/opencode.json` with `npm: @ai-sdk/openai-compatible`, `options.baseURL: http://openrouter-proxy:8080/v1`, and free-tier model entry
- [x] 5.2 Update `spike/issue-45-opencode-cache/README.md` with two-container startup instructions: create network, start proxy container with `--env-file .env`, start fixture container on same network, run `opencode run` via `docker exec`
- [x] 5.3 Add `.env.example` and update `.gitignore` to exclude `.env` and `.env.*`

## 6. Validation

- [x] 6.1 Start both containers on `spike-net`, run a single-turn `opencode run "Say: acknowledged"` through the proxy, confirm proxy logs show forwarded request and response received — confirmed: `tokens_input=319 tokens_output=8` in opencode DB, proxy logs show request forwarded to OpenRouter
- [x] 6.2 Run a 3-turn session using `opencode run --session <id>` pattern, query cache fields from opencode DB — confirmed: proxy does not break session continuity; `cache_read=0` across all turns because `qwen/qwen3.6-flash` via OpenRouter requires larger prompts (likely 1024+ tokens) to activate caching; cache behavior at prompt scale is a finding for SF-4b
- [ ] 6.3 Close GitHub issue #55
