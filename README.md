# Retrace — Where did I see that?

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8.svg)](https://go.dev)
[![Local-first](https://img.shields.io/badge/Local--first-100%25-brightgreen.svg)](#privacy)
[![No telemetry](https://img.shields.io/badge/Telemetry-none-red.svg)](#privacy)

> You remember reading it. You don't remember where.
> A link in an email. An address in a chat. A recipe in a note.
> **Retrace finds it — one query across everything you own.**

Most of what you know is scattered: Gmail, Telegram, files, browser history, notes.
Each of them has its own search box, and none of them talks to the others.
So you re-read the same thread three times, or give up and ask the person again.

Retrace is a personal search engine that runs entirely on your machine.
It indexes your stuff locally and answers one question: **where did I see that?**

## What it does

```
$ retrace index ~/notes ~/docs
indexed 1,284 files (312 skipped)

$ retrace search "lighthouse Kvirike"
~/notes/trips/georgia.md
  …The lighthouse at Kvirike is worth the climb, best at sunset…
```

- **One query, every source.** Files today; mail, chats and browser history next.
- **Full-text search that understands words**, not just substrings — stemming, ranking, snippets.
- **Instant.** A local SQLite FTS5 index; no server round-trip, no loading spinner.
- **Private by construction.** There is no account, no cloud, no telemetry. The index lives in `~/.local/share/retrace/` and never leaves your disk. Retrace makes zero network requests, ever.

## Quick start

```bash
git clone https://github.com/ranas-mukminov/retrace
cd retrace
go build -o retrace ./cmd/retrace

./retrace index ~/notes ~/Documents   # index your text files
./retrace search "plumber phone"      # find it
./retrace status                      # how much is indexed
```

Requirements: Go 1.26+. No C compiler needed — the SQLite driver is pure Go,
so the same binary builds and runs on every Linux distribution.

## What works now (v0.1)

- [x] Local full-text index (SQLite FTS5, porter stemming, bm25 ranking)
- [x] Incremental indexing — unchanged files are skipped
- [x] `index`, `search`, `status` CLI
- [x] Snippets around every match
- [x] Zero network access, zero telemetry

## Roadmap

- [ ] **Connectors** — Gmail (read-only, OAuth), Telegram export, browser history (Firefox/Chrome), notes apps
- [ ] **Semantic search** — local embeddings, so "that article about sleep" finds the article that never says "sleep"
- [ ] **Attachments** — PDFs, office docs, images (OCR)
- [ ] **GUI + global hotkey** — press a key anywhere, type, jump to the source
- [ ] **Watch mode** — the index follows your files as they change
- [ ] **Packages** — AppImage, Flatpak, .deb, .rpm

The roadmap is the plan, not a promise — checked items are the only ones that exist.

## Privacy

This is the whole point, so it is stated plainly:

- **Local-first.** All indexing and searching happens on your machine.
- **No account.** Nothing to sign up for, nothing to log in to.
- **No telemetry.** No analytics, no crash reports, no phone-home of any kind.
- **No network.** Retrace never opens a socket. You can verify that: the code has no HTTP client.
- **Your data stays yours.** The index is a file on your disk. Delete it and it is gone.

Future connectors (Gmail, Telegram) will ask for the minimum read-only access
they need, and will say so in plain language before you grant it.

## Architecture

```
sources ──► connectors ──► index (SQLite FTS5) ──► query ──► CLI / GUI
   files         files           local file            bm25
   (mail)        (mail…)         ~/.local/share        snippets
   (chats)                       retrace/retrace.db
```

One small core, many thin connectors. Adding a source means writing one
connector that feeds documents into the index — nothing else changes.
See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Contributing

Small, focused PRs. One change per PR, with a test where it makes sense.
See [CONTRIBUTING.md](CONTRIBUTING.md). Issues and ideas are welcome —
especially "I searched for X and couldn't find it".

## License

MIT. See [LICENSE](LICENSE).
