package exporter

import (
	"strings"
	"testing"
	"time"

	"agentdata/internal/model"
)

func TestJSONLWritesOneRecordPerMessage(t *testing.T) {
	session := sampleSession()
	var out strings.Builder

	err := JSONL(&out, []model.Session{session})
	if err != nil {
		t.Fatalf("JSONL returned error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, `"session_id":"s1"`) {
		t.Fatalf("JSONL output missing session_id: %s", got)
	}
	if !strings.Contains(got, `"text":"Fix the build"`) {
		t.Fatalf("JSONL output missing message text: %s", got)
	}
	if strings.Count(strings.TrimSpace(got), "\n") != 1 {
		t.Fatalf("JSONL output = %q, want two lines", got)
	}
}

func TestMarkdownGroupsMessagesBySession(t *testing.T) {
	session := sampleSession()
	var out strings.Builder

	err := Markdown(&out, []model.Session{session})
	if err != nil {
		t.Fatalf("Markdown returned error: %v", err)
	}

	got := out.String()
	for _, want := range []string{"# s1", "Source: codex", "## user", "Fix the build", "## assistant", "Use go test ./..."} {
		if !strings.Contains(got, want) {
			t.Fatalf("Markdown output missing %q: %s", want, got)
		}
	}
}

func sampleSession() model.Session {
	return model.Session{
		ID:      "s1",
		Source:  "codex",
		Project: "D:\\go_project",
		Messages: []model.Message{
			{Role: "user", Text: "Fix the build", Timestamp: time.Date(2026, 5, 11, 1, 0, 0, 0, time.UTC)},
			{Role: "assistant", Text: "Use go test ./...", Timestamp: time.Date(2026, 5, 11, 1, 1, 0, 0, time.UTC)},
		},
	}
}
