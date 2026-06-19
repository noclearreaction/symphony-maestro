# openrouter-proxy — Docker usage

A minimal Go HTTP proxy between opencode and OpenRouter. Forwards `POST /v1/chat/completions` and passes through SSE streaming. Injects `Authorization: Bearer` from the `OPENROUTER_API_KEY` environment variable.

## Build

```bash
docker build -t openrouter-proxy spike/issue-45-opencode-cache/proxy/
```

Run from repo root.

## Run

```bash
docker run -d --name openrouter-proxy --network spike-net \
  --env-file .env \
  -v /tmp/proxy-logs:/logs \
  openrouter-proxy
```

Each forwarded request is appended as one JSON line to a session-keyed NDJSON file in the log directory. The session key is derived from the SHA-256 hash of the first 512 bytes of `messages[0].content` — requests sharing the same stable system prompt land in the same file. A 3-turn opencode session produces 2 files: one for the agent turns (e.g. `d46d701d.ndjson`) and one for the title-generation request (`12effdd3.ndjson`).

Each line in the NDJSON file is a self-contained JSON object:
```json
{"timestamp":"...","turn":1,"request":{...},"response":{...}}
```

Inspect with:
```bash
cat /tmp/proxy-logs/<session-key>.ndjson | while read line; do echo "$line" | python3 -m json.tool; done
```

The container listens on port 8080. It must share a Docker network with the fixture container so opencode can reach it at `http://openrouter-proxy:8080`.

## Source

- `main.go` — single-file Go proxy, stdlib only, no external dependencies
- `go.mod` — module definition
- `Dockerfile` — multi-stage build: `golang:1.26` builder → `alpine:3.21` runtime

## Environment

| Variable | Required | Description |
|---|---|---|
| `OPENROUTER_API_KEY` | Yes | OpenRouter API key. Without it the proxy starts but all requests will be rejected by OpenRouter. |
| `PORT` | No | Port to listen on. Defaults to `8080`. |
| `LOG_DIR` | No | Directory to write per-request JSON log files. Defaults to `/logs`. Created on startup if absent. Mount a host directory here to collect logs. |

## Routes

| Route | Behaviour |
|---|---|
| `POST /v1/chat/completions` | Forwarded to `https://openrouter.ai/api/v1/chat/completions` |
| All other routes | `404 Not Found` |
