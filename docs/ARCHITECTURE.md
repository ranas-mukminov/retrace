# Architecture

Retrace is one small core with many thin connectors.

## Components

- **Connectors** read one source and produce documents (path + content).
  v0.1 has one: the file connector (`internal/index`), which walks
  directories and reads text files.
- **Index** is a SQLite database with an FTS5 virtual table, stored at
  `~/.local/share/retrace/retrace.db` (mode 0700). Keyword search uses
  porter stemming and bm25 ranking; `snippet()` produces the context lines.
- **Query** is the CLI today (`cmd/retrace`); a GUI and a global hotkey
  will sit on the same index later.

## Data flow

```
sources ──► connectors ──► index (SQLite FTS5) ──► query ──► CLI / GUI
```

Adding a source means writing one connector that feeds documents into the
index. Nothing else changes: ranking, snippets and storage are shared.

## Privacy invariants

- No network access anywhere in the codebase. There is no HTTP client.
- No telemetry, no account, no crash reporting.
- The index is a local file. Deleting it removes everything Retrace knows.
- Connectors request the minimum read-only access they need and document it.
