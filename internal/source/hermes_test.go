package source

import (
	"strings"
	"testing"
	"time"
)

func TestParseHermesTranscriptNormalizesMessages(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		`{"session_id":"hermes-session-1","project":"D:\\go_project","timestamp":"2026-05-12T01:02:00Z","role":"user","content":"Find deploy notes"}`,
		`{"session_id":"hermes-session-1","project":"D:\\go_project","timestamp":"2026-05-12T01:03:00Z","role":"assistant","content":"Found them in the release session."}`,
	}, "\n"))

	session, err := ParseHermesTranscript(input, "hermes-session-1.jsonl")
	if err != nil {
		t.Fatalf("ParseHermesTranscript returned error: %v", err)
	}

	if session.ID != "hermes-session-1" {
		t.Fatalf("session.ID = %q, want hermes-session-1", session.ID)
	}
	if session.Source != "hermes" {
		t.Fatalf("session.Source = %q, want hermes", session.Source)
	}
	if session.Project != `D:\go_project` {
		t.Fatalf("session.Project = %q, want D:\\go_project", session.Project)
	}
	if len(session.Messages) != 2 {
		t.Fatalf("len(session.Messages) = %d, want 2", len(session.Messages))
	}
	if session.Messages[0].Role != "user" || session.Messages[0].Text != "Find deploy notes" {
		t.Fatalf("first message = %+v, want normalized user message", session.Messages[0])
	}

	wantTime := time.Date(2026, 5, 12, 1, 2, 0, 0, time.UTC)
	if !session.Messages[0].Timestamp.Equal(wantTime) {
		t.Fatalf("timestamp = %s, want %s", session.Messages[0].Timestamp, wantTime)
	}
}

func TestParseHermesTranscriptAcceptsMessageObject(t *testing.T) {
	input := strings.NewReader(`{"sessionId":"s1","message":{"role":"assistant","content":[{"type":"text","text":"hello"}]}}`)

	session, err := ParseHermesTranscript(input, "s1.jsonl")
	if err != nil {
		t.Fatalf("ParseHermesTranscript returned error: %v", err)
	}
	if len(session.Messages) != 1 || session.Messages[0].Text != "hello" {
		t.Fatalf("messages = %+v, want text block content", session.Messages)
	}
}
