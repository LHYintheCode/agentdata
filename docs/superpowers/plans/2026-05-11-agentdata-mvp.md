# Agentdata MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create a tested local-first Go CLI for scanning, searching, and exporting normalized agent chat records.

**Architecture:** Keep vendor formats behind source adapters. Use canonical model types across search, export, and CLI layers. Start with explicit-path JSONL ingestion before automatic source discovery.

**Tech Stack:** Go 1.26, standard library, `go test ./...`, GitHub CLI for remote repository creation once authenticated.

---

### Task 1: Project Skeleton

**Files:**
- Create: `go.mod`
- Create: `README.md`
- Create: `.gitignore`

- [ ] Initialize Go module.
- [ ] Add README with product positioning and local-only privacy defaults.
- [ ] Add gitignore for built binaries and temporary data.
- [ ] Run `go test ./...`.
- [ ] Commit skeleton.

### Task 2: Canonical Model and JSONL Source Parser

**Files:**
- Create: `internal/model/model.go`
- Create: `internal/source/jsonl.go`
- Create: `internal/source/jsonl_test.go`

- [ ] Write failing tests for parsing role/content/timestamp/project/session fields from JSONL.
- [ ] Implement canonical types and parser.
- [ ] Run `go test ./internal/source`.
- [ ] Commit parser.

### Task 3: Local Search

**Files:**
- Create: `internal/search/search.go`
- Create: `internal/search/search_test.go`

- [ ] Write failing tests for case-insensitive message search.
- [ ] Implement search over canonical sessions.
- [ ] Run `go test ./internal/search`.
- [ ] Commit search.

### Task 4: Exporters

**Files:**
- Create: `internal/exporter/exporter.go`
- Create: `internal/exporter/exporter_test.go`

- [ ] Write failing tests for JSONL and Markdown export.
- [ ] Implement exporters.
- [ ] Run `go test ./internal/exporter`.
- [ ] Commit exporters.

### Task 5: CLI Commands

**Files:**
- Create: `cmd/agentdata/main.go`
- Create: `cmd/agentdata/main_test.go`

- [ ] Write failing CLI tests for `version`, `scan`, `search`, and `export`.
- [ ] Implement CLI using the standard `flag` package.
- [ ] Run `go test ./...`.
- [ ] Commit CLI.

### Task 6: GitHub Publication

**Files:**
- Modify: `README.md`

- [ ] Verify `gh auth status`.
- [ ] Create public repository `agentdata`.
- [ ] Add remote origin.
- [ ] Push all local commits.
- [ ] Add a clear README note that the project is local-only and early-stage.
