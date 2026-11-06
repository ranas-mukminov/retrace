# Contributing to Retrace

Issues and ideas are welcome — especially "I searched for X and couldn't find it".

## Pull requests

- Small and focused: one change per PR.
- Add a test where the change has behavior worth pinning down.
- Run `gofmt`, `go vet ./...` and `go test ./...` before pushing.
- Keep the privacy promise: no network access, no telemetry, no account.
  A PR that phones home will not be merged.

## Adding a connector

A connector reads one source (mail, chat export, browser history) and feeds
documents into the index. It should:

1. Ask for the minimum access it needs, and say so in plain language.
2. Never send anything anywhere — documents go to the local index only.
3. Come with a short doc: what it reads, what it stores, how to revoke it.
