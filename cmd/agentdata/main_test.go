package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := Run([]string{"version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "agentdata") {
		t.Fatalf("stdout = %q, want version text", stdout.String())
	}
}

func TestRunScanCountsSessionsAndMessages(t *testing.T) {
	dir := writeSampleJSONL(t)
	var stdout, stderr bytes.Buffer

	code := Run([]string{"scan", "--path", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "sessions=1 messages=2") {
		t.Fatalf("stdout = %q, want scan counts", stdout.String())
	}
}

func TestRunScanCodexSourceCountsSessionsAndMessages(t *testing.T) {
	dir := writeSampleCodexRollout(t)
	var stdout, stderr bytes.Buffer

	code := Run([]string{"scan", "--source", "codex", "--path", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "sessions=1 messages=2") {
		t.Fatalf("stdout = %q, want scan counts", stdout.String())
	}
}

func TestRunScanClaudeSourceCountsSessionsAndMessages(t *testing.T) {
	dir := writeSampleClaudeTranscript(t)
	var stdout, stderr bytes.Buffer

	code := Run([]string{"scan", "--source", "claude", "--path", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "sessions=1 messages=2") {
		t.Fatalf("stdout = %q, want scan counts", stdout.String())
	}
}

func TestRunScanAllSourceCountsCodexAndClaude(t *testing.T) {
	codexDir := writeSampleCodexRollout(t)
	claudeDir := writeSampleClaudeTranscript(t)
	var stdout, stderr bytes.Buffer

	code := Run([]string{"scan", "--source", "all", "--path", "codex=" + codexDir + ",claude=" + claudeDir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "sessions=2 messages=4") {
		t.Fatalf("stdout = %q, want combined scan counts", stdout.String())
	}
}

func TestRunSearchPrintsMatches(t *testing.T) {
	dir := writeSampleJSONL(t)
	var stdout, stderr bytes.Buffer

	code := Run([]string{"search", "--path", dir, "deploy"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	got := stdout.String()
	if !strings.Contains(got, "s1") || !strings.Contains(got, "Deploy the CLI") {
		t.Fatalf("stdout = %q, want matching session and text", got)
	}
}

func TestRunSearchCodexSourcePrintsMatches(t *testing.T) {
	dir := writeSampleCodexRollout(t)
	var stdout, stderr bytes.Buffer

	code := Run([]string{"search", "--source", "codex", "--path", dir, "deploy"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	got := stdout.String()
	if !strings.Contains(got, "codex-session-1") || !strings.Contains(got, "Deploy the CLI") {
		t.Fatalf("stdout = %q, want matching codex session and text", got)
	}
}

func TestRunSearchClaudeSourcePrintsMatches(t *testing.T) {
	dir := writeSampleClaudeTranscript(t)
	var stdout, stderr bytes.Buffer

	code := Run([]string{"search", "--source", "claude", "--path", dir, "deploy"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	got := stdout.String()
	if !strings.Contains(got, "claude-session-1") || !strings.Contains(got, "Find my deploy notes") {
		t.Fatalf("stdout = %q, want matching claude session and text", got)
	}
}

func TestRunSearchAllSourcePrintsMatchesAcrossSources(t *testing.T) {
	codexDir := writeSampleCodexRollout(t)
	claudeDir := writeSampleClaudeTranscript(t)
	var stdout, stderr bytes.Buffer

	code := Run([]string{"search", "--source", "all", "--path", "codex=" + codexDir + ",claude=" + claudeDir, "deploy"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	got := stdout.String()
	if !strings.Contains(got, "codex-session-1") || !strings.Contains(got, "claude-session-1") {
		t.Fatalf("stdout = %q, want matches from codex and claude", got)
	}
}

func TestRunExportMarkdown(t *testing.T) {
	dir := writeSampleJSONL(t)
	var stdout, stderr bytes.Buffer

	code := Run([]string{"export", "--path", dir, "--format", "markdown"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "# s1") {
		t.Fatalf("stdout = %q, want markdown session", stdout.String())
	}
}

func TestRunExportWritesOutputFile(t *testing.T) {
	dir := writeSampleJSONL(t)
	outPath := filepath.Join(t.TempDir(), "history.md")
	var stdout, stderr bytes.Buffer

	code := Run([]string{"export", "--path", dir, "--format", "markdown", "--out", outPath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want no stdout when --out is used", stdout.String())
	}
	contents, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read out file: %v", err)
	}
	if !strings.Contains(string(contents), "# s1") {
		t.Fatalf("out file = %q, want markdown session", string(contents))
	}
}

func TestResolveSourcePathDefaultsCodexSessions(t *testing.T) {
	path, err := resolveSourcePath("codex", "")
	if err != nil {
		t.Fatalf("resolveSourcePath returned error: %v", err)
	}
	if !strings.HasSuffix(path, filepath.Join(".codex", "sessions")) {
		t.Fatalf("path = %q, want .codex sessions path", path)
	}
}

func TestResolveSourcePathDefaultsClaudeProjects(t *testing.T) {
	path, err := resolveSourcePath("claude", "")
	if err != nil {
		t.Fatalf("resolveSourcePath returned error: %v", err)
	}
	if !strings.HasSuffix(path, filepath.Join(".claude", "projects")) {
		t.Fatalf("path = %q, want .claude projects path", path)
	}
}

func TestResolveSourcePathRequiresPathForJSONL(t *testing.T) {
	_, err := resolveSourcePath("jsonl", "")
	if err == nil {
		t.Fatal("resolveSourcePath returned nil error for missing JSONL path")
	}
}

func writeSampleJSONL(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	data := strings.Join([]string{
		`{"source":"codex","project":"D:\\go_project","session_id":"s1","timestamp":"2026-05-11T01:02:03Z","role":"user","content":"Deploy the CLI"}`,
		`{"source":"codex","project":"D:\\go_project","session_id":"s1","timestamp":"2026-05-11T01:03:04Z","role":"assistant","content":"Run go test ./..."}`,
	}, "\n")
	if err := os.WriteFile(filepath.Join(dir, "sample.jsonl"), []byte(data), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	return dir
}

func writeSampleCodexRollout(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	data := strings.Join([]string{
		`{"timestamp":"2026-05-11T01:00:00Z","type":"session_meta","payload":{"id":"codex-session-1","cwd":"D:\\go_project"}}`,
		`{"timestamp":"2026-05-11T01:02:00Z","type":"response_item","payload":{"role":"user","content":[{"type":"input_text","text":"Deploy the CLI"}]}}`,
		`{"timestamp":"2026-05-11T01:03:00Z","type":"response_item","payload":{"role":"assistant","content":[{"type":"output_text","text":"Run go test ./..."}]}}`,
	}, "\n")
	if err := os.WriteFile(filepath.Join(dir, "rollout-2026-05-11T01-00-00-codex-session-1.jsonl"), []byte(data), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	return dir
}

func writeSampleClaudeTranscript(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	data := strings.Join([]string{
		`{"type":"user","sessionId":"claude-session-1","cwd":"D:\\go_project","timestamp":"2026-05-11T01:02:00.000Z","message":{"role":"user","content":"Find my deploy notes"}}`,
		`{"type":"assistant","sessionId":"claude-session-1","cwd":"D:\\go_project","timestamp":"2026-05-11T01:03:00.000Z","message":{"role":"assistant","content":[{"type":"text","text":"The deploy notes are in session 42."}]}}`,
	}, "\n")
	if err := os.WriteFile(filepath.Join(dir, "claude-session-1.jsonl"), []byte(data), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	return dir
}
