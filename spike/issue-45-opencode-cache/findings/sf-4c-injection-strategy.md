# SF-4c: Injection Strategy Findings

## Context

The proxy sits between opencode and OpenRouter. SF-4c injects `cache_control` markers and potentially dynamic state into forwarded requests. This document records what was learned about the message structure, caching mechanics, and injection strategy options before implementation.

## Message Array Structure

opencode sends a standard OpenAI-compatible messages array. Every message has exactly two keys: `role` and `content` (always a plain string — never content-parts arrays from opencode's side).

```
messages[0]   role=system     stable system prompt (INSTRUCTIONS.md)
messages[1]   role=user       first user turn
messages[2]   role=assistant  model reply to turn 1
messages[3]   role=user       second user turn
...
messages[-1]  role=user       current user turn (always the last entry)
```

History grows by 2 entries per turn. `messages[-1]` is always the current user input.

## How Caching Works (KV Prefix Cache)

The model server caches KV activations for token prefixes. When a new request arrives, it finds the longest matching prefix in its cache and reuses those activations, charging a reduced rate for cached tokens.

**Key properties:**

- The cache is **read-only** per request — it is a compute shortcut, not shared mutable state
- Two sessions with identical `messages` arrays share the same cache hits but produce independent outputs
- **No cross-hearing**: each request's context window is exactly its own `messages[]`, nothing more
- With `temperature=0` and identical messages, outputs are deterministic — but the model is unaware of other requests

## Cache Sharing Across Sessions

For models with automatic/implicit caching (e.g. Gemini 2.5 Flash), any two requests sharing a common leading token sequence share cache benefit. This extends progressively:

- Turn 1 of any new session hits the system prompt cache immediately — no warmup cost
- If conversations stay identical (e.g. deterministic agent tasks), the shared prefix grows each turn
- A divergence in any assistant response breaks prefix matching from that point forward

**Implication for this project:** The intended use pattern — restarting the same agent model repeatedly with the same system prompt — is the **ideal caching scenario**. Turn 1 of every new session is already a cache hit. There is no first-session warmup tax.

## Injection Strategy Options

There are two places the proxy could inject dynamic state:

### Option 1: Append to `messages[0]` (system prompt)

For models with explicit caching (e.g. Anthropic Claude), `cache_control` markers can be placed on individual content blocks within a message's content array. This allows:

```json
{
  "role": "system",
  "content": [
    {"type": "text", "text": "<stable instructions>", "cache_control": {"type": "ephemeral"}},
    {"type": "text", "text": "<dynamic state>"}
  ]
}
```

The cached block is preserved; only the dynamic suffix is recomputed fresh. This works for Anthropic. **Gemini's behavior through OpenRouter is unverified** — Gemini uses a separate `cachedContent` mechanism natively; whether OpenRouter's translation supports a mid-message cache boundary is unknown and requires empirical testing.

**Risk:** Any change to content before the cache marker is a full cache eviction.

### Option 2: Prepend to `messages[-1]` (current user turn)

Append a state header to the current user message before forwarding:

```
[proxy-state: turn=3, cached_tokens=1164, cost_so_far=$0.0003]
What is 10-3?
```

`messages[-1]` is always past the cache boundary. The stable system prompt cache is **never touched**. Works identically for all models and caching mechanisms.

**Tradeoff:** Injected content in `messages[-1]` becomes part of conversation history on the next turn (`messages[-2]` will contain it). If the injection is non-deterministic or session-specific, it will prevent cross-session cache sharing at the history level. If it must be deterministic for shared cache benefit, inject only values that are identical across sessions.

## Decision

**SF-4c will use Option 2 only: prepend dynamic state to `messages[-1]`.**

Option 1 (cache_control markers, content-array transformation) is **abandoned** for the current infrastructure. Rationale:

1. Works for all models without model-specific content-array transformation
2. Never invalidates the system prompt cache — `messages[0]` stays frozen
3. The use pattern (frequent restarts of the same agent model) makes system prompt cache reuse more valuable than history cache sharing
4. Brief injections (a few tokens of state data) are negligible in cost compared to the round-trip savings from avoiding extra tool-only turns
5. The only adversarial case — two sessions with identical `messages[-1]` content that would otherwise share deep history cache — does not apply once session-specific state is being injected

**Keep `messages[0]` frozen.** Treat a change to `INSTRUCTIONS.md` as a cache eviction event for all running instances. Versioned, additive changes to instructions are cheaper than in-place edits.

**Injection format:** a brief structured prefix prepended to the user message content, e.g.:
```
[state: turn=3, cached_tokens=1164]
What is 10-3?
```

The prefix should be kept minimal. Verbose injection defeats the purpose.

## Injection Selection: Magic Marker Design (Resolved)

The proxy scans `messages[0].content` for a `## Runtime state injection` section. If found, injection is active for that agent. If absent, the request is forwarded unchanged.

The system prompt (`INSTRUCTIONS.md` or `prompt.md`) contains a dedicated marker block. The surrounding prose is for the model; the marker block is for the proxy:

```markdown
## Runtime state injection

Every message will be prepended with runtime information listed below.
The data is current as of the moment the message was sent and can be relied upon.
Use it to avoid asking for information the system already provides.

<!-- proxy:inject
git_status: current branch, uncommitted changes, and last commit summary
unittests: output of `go test ./...` — pass/fail status and elapsed time
-->
```

The proxy searches `messages[0].content` for `<!-- proxy:inject` and extracts everything up to `-->`. Each line inside the block is a `key: description` pair. The proxy populates each key with current state data and prepends the result to `messages[-1].content`:

```
runtime_state:
  git_status: "branch: feature/issue-55, 3 uncommitted changes"
  unittests: "ok github.com/... (0.14s)"

<original user message>
```

**Why this design works:**

- The proxy marker (`<!-- proxy:inject ... -->`) is the only frozen string — the section heading and surrounding prose are free to change without cache consequence
- The model reads the comment block as natural language (HTML comments are visible in raw text passed to the model) — it sees the description alongside the key name and knows what to expect each turn
- The declared keys in the block ARE the proxy configuration — no env vars or separate config file needed
- Title-gen and any agent without the block are forwarded unchanged
- The description in the block serves the model; the key name serves the proxy

**Proxy implementation sketch:**

1. Parse `messages[0].content` — search for `<!-- proxy:inject`
2. If found, extract lines between the marker and `-->` as `key: description` pairs
3. For each key, invoke the registered data source
4. Prepend `runtime_state:\n  key: value\n  ...` to `messages[-1].content`
5. Forward the mutated request

**Constraint:** Only the `<!-- proxy:inject` and `-->` delimiters and the key names within are frozen. Descriptions, heading text, and surrounding prose are free to evolve.

## Cross-Session Log Correlation

The proxy cannot observe opencode's session ID — it is never sent over the wire. The proxy log's session key (derived from `messages[0].content` hash) groups by agent identity (system prompt), not by conversation.

To correlate proxy logs with opencode's SQLite session records, hash `messages[-1].content` **before any mutation** and store it in the log entry as `last_user_hash`. opencode's DB stores the original user message; `sha256(original_content)[:8]` produces the same hash on the other side.

```json
{"timestamp":"...","turn":2,"last_user_hash":"a3f8c1d2","request":{...},"response":{...}}
```

This join key survives proxy mutation of the forwarded request.
