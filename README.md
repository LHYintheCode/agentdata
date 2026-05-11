# Agentdata

Agentdata is a local-first CLI for developers who want to search, normalize, and export their coding-agent chat history.

The project starts with explicit-path JSONL ingestion and a small canonical schema. It is designed for future adapters for Codex, Claude Code, Trae, and other agent tools without tying the core data model to any one vendor format.

## Privacy Defaults

- Local-only by default.
- Read-only scanning.
- No telemetry.
- No cloud sync.
- Exports go to stdout in the initial MVP.

## MVP Commands

```powershell
go run ./cmd/agentdata version
go run ./cmd/agentdata scan --path ./samples
go run ./cmd/agentdata search --path ./samples "deploy"
go run ./cmd/agentdata export --path ./samples --format markdown
```

## Status

Early-stage public CLI scaffold. Marketplace, data sale, consent manifests, redaction, and source auto-discovery are intentionally out of scope for the first implementation slice.
