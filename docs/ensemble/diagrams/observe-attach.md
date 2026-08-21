# Observing and attaching

Observing takes nothing from a running part; attaching takes ownership.
Companion to [README §5](../README.md#5-dispatch-flows).

```mermaid
sequenceDiagram
  autonumber
  participant H as human
  participant E as Ensemble
  participant P as part
  H->>E: observe(auth-refactor)
  activate E
  E-->>H: session/update stream
  Note over P: still playing
  deactivate E
  H->>E: attach(auth-refactor)
  Note over E: only from AwaitingReview or Faulted
  E-->>H: terminal handed over
  Note over E: queued work waits until detach
```
