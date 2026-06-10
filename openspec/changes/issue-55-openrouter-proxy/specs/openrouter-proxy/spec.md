## ADDED Requirements

### Requirement: Proxy accepts OpenAI-compatible chat completion requests
The proxy SHALL accept `POST /v1/chat/completions` requests using the OpenAI chat completions wire format. All other routes SHALL return HTTP 404.

#### Scenario: Valid chat completion request forwarded
- **WHEN** opencode sends `POST /v1/chat/completions` with a valid JSON body
- **THEN** the proxy forwards the request to OpenRouter and returns the response to opencode

#### Scenario: Unknown route rejected
- **WHEN** a request arrives at any path other than `/v1/chat/completions`
- **THEN** the proxy returns HTTP 404

---

### Requirement: Proxy forwards SSE streaming responses without buffering
The proxy SHALL forward Server-Sent Events streaming responses from OpenRouter to the client incrementally, without accumulating the full response body before forwarding.

#### Scenario: Streaming response forwarded in real time
- **WHEN** OpenRouter returns a `text/event-stream` response
- **THEN** each SSE chunk is forwarded to opencode as it arrives, with no buffering delay

#### Scenario: Non-streaming response forwarded correctly
- **WHEN** OpenRouter returns a non-streaming JSON response
- **THEN** the full response body is forwarded to opencode with correct `Content-Type`

---

### Requirement: Proxy optionally sets Authorization header with OpenRouter API key
If `OPENROUTER_API_KEY` is set in the proxy's environment, the proxy SHALL set `Authorization: Bearer ${OPENROUTER_API_KEY}` on all forwarded requests. If the key is not set, the proxy SHALL forward requests without an `Authorization` header. The key SHALL only be present in the proxy's environment, never in opencode's configuration.

#### Scenario: Key present — forwarded with Authorization header
- **WHEN** `OPENROUTER_API_KEY` is set and opencode sends a request
- **THEN** the proxy forwards the request with `Authorization: Bearer ${OPENROUTER_API_KEY}`

#### Scenario: Key absent — forwarded without Authorization header
- **WHEN** `OPENROUTER_API_KEY` is not set and opencode sends a request
- **THEN** the proxy forwards the request without an `Authorization` header and logs a warning at startup

#### Scenario: Missing key does not prevent startup
- **WHEN** the proxy starts and `OPENROUTER_API_KEY` is not set
- **THEN** the proxy starts normally, logs a warning, and continues to serve requests

---

### Requirement: Proxy and fixture run on a shared Docker network
The proxy SHALL be deployable as a Docker container on a user-defined Docker network, reachable by the fixture container using the proxy's container name as hostname.

#### Scenario: opencode connects to proxy by container hostname
- **WHEN** the fixture container runs with `OPENROUTER_PROXY_URL=http://openrouter-proxy:8080`
- **THEN** opencode successfully routes all AI requests through the proxy

#### Scenario: Proxy and fixture started with documented two-container workflow
- **WHEN** a user follows the README startup instructions
- **THEN** both containers start, connect, and a single `opencode run` turn completes successfully through the proxy
