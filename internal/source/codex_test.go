package source

import (
	"strings"
	"testing"
	"time"
)

func TestParseCodexRolloutNormalizesUserAndAssistantMessages(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		`{"timestamp":"2026-05-11T01:00:00Z","type":"session_meta","payload":{"id":"codex-session-1","cwd":"D:\\go_project","source":"vscode"}}`,
		`{"timestamp":"2026-05-11T01:01:00Z","type":"response_item","payload":{"type":"message","role":"developer","content":[{"type":"input_text","text":"internal instructions"}]}}`,
		`{"timestamp":"2026-05-11T01:02:00Z","type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"Find old deploy notes"}]}}`,
		`{"timestamp":"2026-05-11T01:03:00Z","type":"response_item","payload":{"type":"message","role":"assistant","content":[{"type":"output_text","text":"The deploy notes are in session 42."}]}}`,
	}, "\n"))

	session, err := ParseCodexRollout(input, "rollout-2026-05-11T01-00-00-codex-session-1.jsonl")
	if err != nil {
		t.Fatalf("ParseCodexRollout returned error: %v", err)
	}

	if session.ID != "codex-session-1" {
		t.Fatalf("session.ID = %q, want codex-session-1", session.ID)
	}
	if session.Source != "codex" {
		t.Fatalf("session.Source = %q, want codex", session.Source)
	}
	if session.Project != `D:\go_project` {
		t.Fatalf("session.Project = %q, want D:\\go_project", session.Project)
	}
	if len(session.Messages) != 2 {
		t.Fatalf("len(session.Messages) = %d, want 2", len(session.Messages))
	}
	if session.Messages[0].Role != "user" || session.Messages[0].Text != "Find old deploy notes" {
		t.Fatalf("first message = %+v, want normalized user message", session.Messages[0])
	}

	wantTime := time.Date(2026, 5, 11, 1, 2, 0, 0, time.UTC)
	if !session.Messages[0].Timestamp.Equal(wantTime) {
		t.Fatalf("timestamp = %s, want %s", session.Messages[0].Timestamp, wantTime)
	}
}

func TestParseCodexRolloutFallsBackToFilenameID(t *testing.T) {
	input := strings.NewReader(`{"timestamp":"2026-05-11T01:02:00Z","type":"response_item","payload":{"role":"user","content":[{"type":"input_text","text":"hello"}]}}`)

	session, err := ParseCodexRollout(input, `C:\Users\me\.codex\sessions\rollout-abc.jsonl`)
	if err != nil {
		t.Fatalf("ParseCodexRollout returned error: %v", err)
	}

	if session.ID != "rollout-abc" {
		t.Fatalf("session.ID = %q, want rollout-abc", session.ID)
	}
}

func TestParseCodexRolloutAcceptsLongLines(t *testing.T) {
	longText := strings.Repeat("x", 128*1024)
	input := strings.NewReader(`{"timestamp":"2026-05-11T01:02:00Z","type":"response_item","payload":{"role":"user","content":[{"type":"input_text","text":"` + longText + `"}]}}`)

	session, err := ParseCodexRollout(input, "rollout-long.jsonl")
	if err != nil {
		t.Fatalf("ParseCodexRollout returned error: %v", err)
	}
	if len(session.Messages) != 1 {
		t.Fatalf("len(session.Messages) = %d, want 1", len(session.Messages))
	}
}
