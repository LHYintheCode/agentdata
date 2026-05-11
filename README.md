# Agentdata

[Chinese](README.zh-CN.md)

Agentdata is a local-first CLI for searching, normalizing, and exporting coding-agent chat history.

It is built for developers who use tools like Codex, Claude Code, Trae, Cursor, or Windsurf and want a portable local archive of their AI work history.

## Why

AI conversations are becoming part of the developer workflow: debugging sessions, design decisions, tool calls, failed attempts, and final fixes often live only inside one agent product.

Agentdata gives that history a local, inspectable format so you can:

- Search past agent conversations.
- Export sessions as JSONL or Markdown.
- Normalize records from different tools into one schema.
- Keep sensitive data on your own machine by default.
- Build future consent, redaction, and packaging workflows on top of a clear local format.

## Status

Agentdata is early-stage. The current version supports explicit-path JSONL input. Automatic discovery for Codex, Claude Code, Trae, Cursor, Windsurf, and other tools will be added after more real-machine validation.

Current capabilities:

- JSONL parsing
- Codex rollout parsing from `~/.codex/sessions`
- Claude Code transcript parsing from `~/.claude/projects`
- Canonical session/message model
- Local message search
- JSONL export
- Markdown export
- Combined Codex + Claude scanning with `--source all`
- File output with `--out`
- Basic CLI test coverage

Not included yet:

- Automatic source discovery for non-Codex/Claude tools
- SQLite/FTS indexing
- Secret redaction
- Consent manifests
- Cloud sync
- Marketplace or resale workflows

## Install

From source:

```bash
go install github.com/LHYintheCode/agentdata/cmd/agentdata@latest
```

For local development:

```bash
git clone https://github.com/LHYintheCode/agentdata.git
cd agentdata
go test ./...
go run ./cmd/agentdata version
```

## Quick Start

Create a JSONL file:

```jsonl
{"source":"codex","project":"/path/to/project","session_id":"s1","timestamp":"2026-05-11T01:02:03Z","role":"user","content":"Deploy the CLI"}
{"source":"codex","project":"/path/to/project","session_id":"s1","timestamp":"2026-05-11T01:03:04Z","role":"assistant","content":"Run go test ./..."}
```

Scan it:

```bash
agentdata scan --path ./samples
```

Search it:

```bash
agentdata search --path ./samples "deploy"
```

Export it:

```bash
agentdata export --path ./samples --format markdown
agentdata export --path ./samples --format jsonl
```

Scan local Codex sessions:

```bash
agentdata scan --source codex
agentdata search --source codex "deploy"
agentdata export --source codex --format markdown > codex-history.md
```

Scan local Claude Code sessions:

```bash
agentdata scan --source claude
agentdata search --source claude "deploy"
agentdata export --source claude --format markdown > claude-history.md
```

Scan Codex and Claude together:

```bash
agentdata scan --source all
agentdata search --source all "deploy"
agentdata export --source all --format markdown --out agent-history.md
```

## Commands

```text
agentdata version
agentdata scan --path <file-or-directory>
agentdata scan --source codex [--path <codex-sessions-directory>]
agentdata scan --source claude [--path <claude-projects-directory>]
agentdata scan --source all [--path codex=<dir>,claude=<dir>]
agentdata search --path <file-or-directory> <query>
agentdata search --source codex [--path <codex-sessions-directory>] <query>
agentdata search --source claude [--path <claude-projects-directory>] <query>
agentdata search --source all [--path codex=<dir>,claude=<dir>] <query>
agentdata export --path <file-or-directory> --format jsonl|markdown
agentdata export --source codex [--path <codex-sessions-directory>] --format jsonl|markdown
agentdata export --source claude [--path <claude-projects-directory>] --format jsonl|markdown
agentdata export --source all [--path codex=<dir>,claude=<dir>] --format jsonl|markdown [--out <file>]
```

## Data Model

Agentdata normalizes source records into sessions and messages:

```json
{
  "id": "s1",
  "source": "codex",
  "project": "/path/to/project",
  "messages": [
    {
      "role": "user",
      "text": "Deploy the CLI",
      "timestamp": "2026-05-11T01:02:03Z"
    }
  ]
}
```

Source adapters should translate vendor-specific formats into this model. The core search and export layers should not depend on one agent vendor's internal storage format.

## Privacy

Agentdata is designed around local ownership:

- It reads local files only when you pass a path.
- It does not upload data.
- It does not collect telemetry.
- It does not modify source chat logs.
- Exports are written to stdout in the current MVP.

Future sharing workflows should require explicit consent, redaction, and a machine-readable manifest before any data leaves the user's machine.

## Roadmap

- Source adapters for Codex and Claude Code
- Trae/Cursor/Windsurf local storage investigation
- SQLite FTS indexing
- Redaction rules for secrets and personal data
- Export packages with manifest files
- MCP adapter for agent hosts that cannot call the CLI directly

## Development

```bash
go test ./...
go run ./cmd/agentdata version
```

## License

License has not been selected yet.
