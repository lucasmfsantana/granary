# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

Granary exports Granola meeting notes from its local cache to markdown files. It runs as a macOS LaunchAgent every 2 hours, preserving transcripts that Granola may purge from its cache.

## Commands

```bash
# Build
go build -v ./...

# Test all
go test ./...

# Test single package
go test ./exporter/...

# Run export
go run main.go run

# Install/manage LaunchAgent (macOS only)
go run main.go install
go run main.go status
go run main.go uninstall
```

CI runs `go build ./...` and `go test ./...` on every push/PR (see [.github/workflows/go-tests.yml](.github/workflows/go-tests.yml)).

## Architecture

**Entry point:** [main.go](main.go) — Cobra CLI with 6 commands (run, install, uninstall, status, version, help).

**Data flow:**
```
Granola cache-vN.json
  → exporter/cache.go    — find latest cache file, parse dual formats (v5 string / v6+ object)
  → exporter/document.go — CacheState{Documents, SharedDocuments, Transcripts}
  → exporter/exporter.go — filter exportable docs, build stable filename map, compare + write
  → exporter/extractor.go — if cache lacks transcripts, recover them from existing .md files
  → exporter/formatter.go — render title/date/notes/transcript to markdown
  → ~/.local/share/granola-transcripts/*.md
```

**`exporter/` package** is the core — all logic except CLI commands and LaunchAgent management lives here.

**`service/`** handles macOS LaunchAgent plist generation, `launchctl bootstrap/bootout`, and status checks.

## Key Design Details

- **Dual cache format:** v5 stores nested JSON as a string; v6+ stores it as an object. `cache.go` handles both.
- **Filename collisions:** Same title+date gets an 8-char document ID suffix appended.
- **Idempotent writes:** Files are byte-compared before writing; unchanged files are skipped.
- **Transcript recovery:** `extractor.go` parses `**Speaker:** text` lines from existing markdown to restore transcripts when the cache no longer has them.
- **Speaker mapping:** `"microphone"` → `"Me"`, `"system"` → `"Them"` in `formatter.go`.
