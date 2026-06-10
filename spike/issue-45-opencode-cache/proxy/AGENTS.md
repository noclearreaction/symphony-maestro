# openrouter-proxy — Docker usage

A minimal Go HTTP proxy between opencode and OpenRouter. Forwards `POST /v1/chat/completions` and passes through SSE streaming. Injects `Authorization: Bearer` from the `OPENROUTER_API_KEY` environment variable.

## Build

```bash
docker build -t openrouter-proxy spike/issue-45-opencode-cache/proxy/
```

Run from repo root.

## Run

```bash
docker run -d --name openrouter-proxy --network spike-net --env-file .env openrouter-proxy
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

## Routes

| Route | Behaviour |
|---|---|
| `POST /v1/chat/completions` | Forwarded to `https://openrouter.ai/api/v1/chat/completions` |
| All other routes | `404 Not Found` |
