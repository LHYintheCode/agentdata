# Agentdata Public CLI Design

## Goal

Build a local-first Go CLI that lets developers scan, normalize, search, and export chat records from coding agents such as Claude Code, Codex, and Trae.

## Product Scope

The first public version focuses on local ownership and portability:

- Scan user-approved local paths in read-only mode.
- Normalize discovered chat records into one schema.
- Search normalized messages locally.
- Export selected records as JSONL or Markdown.
- Keep source adapters isolated so unstable vendor formats do not leak into the core model.

The first version does not include a marketplace, cloud sync, user accounts, or paid data exchange. Those require separate consent, privacy, licensing, and compliance work.

## Architecture

The CLI is split into small packages:

- `internal/model`: canonical `Session` and `Message` types.
- `internal/source`: source adapters and JSONL parsing helpers.
- `internal/search`: local in-memory search for the first MVP.
- `internal/exporter`: JSONL and Markdown exporters.
- `cmd/agentdata`: command parsing and user-facing CLI.

Source adapters return canonical sessions. The rest of the system never reads vendor-specific fields directly.

## Commands

- `agentdata scan --path <dir> --format jsonl`
- `agentdata search --path <dir> <query>`
- `agentdata export --path <dir> --format jsonl|markdown`
- `agentdata version`

The MVP accepts explicit paths first. Automatic discovery of `~/.claude`, `~/.codex`, and Trae storage paths comes later after real-machine validation.

## Privacy Defaults

- All commands are local-only.
- Scanning is read-only.
- No telemetry is collected.
- Export commands only write to stdout in the MVP.
- Future packaging must include consent and redaction manifests before any data sharing workflow.

## Testing

Use Go tests for parser behavior, search behavior, exporters, and CLI smoke coverage. New behavior starts with failing tests before implementation.
