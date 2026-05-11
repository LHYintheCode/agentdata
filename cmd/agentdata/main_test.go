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
