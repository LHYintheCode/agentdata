package source

import (
	"strings"
	"testing"
	"time"
)

func TestParseClaudeTranscriptNormalizesUserAndAssistantMessages(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		`{"type":"permission-mode","sessionId":"claude-session-1"}`,
		`{"type":"user","sessionId":"claude-session-1","cwd":"D:\\go_project","timestamp":"2026-05-11T01:02:00.000Z","message":{"role":"user","content":"Find my deploy notes"}}`,
		`{"type":"assistant","sessionId":"claude-session-1","cwd":"D:\\go_project","timestamp":"2026-05-11T01:03:00.000Z","message":{"role":"assistant","content":[{"type":"text","text":"The deploy notes are in session 42."},{"type":"tool_use","name":"Read"}]}}`,
	}, "\n"))

	session, err := ParseClaudeTranscript(input, "claude-session-1.jsonl")
	if err != nil {
		t.Fatalf("ParseClaudeTranscript returned error: %v", err)
	}

	if session.ID != "claude-session-1" {
		t.Fatalf("session.ID = %q, want claude-session-1", session.ID)
	}
	if session.Source != "claude" {
		t.Fatalf("session.Source = %q, want claude", session.Source)
	}
	if session.Project != `D:\go_project` {
		t.Fatalf("session.Project = %q, want D:\\go_project", session.Project)
	}
	if len(session.Messages) != 2 {
		t.Fatalf("len(session.Messages) = %d, want 2", len(session.Messages))
	}
	if session.Messages[0].Role != "user" || session.Messages[0].Text != "Find my deploy notes" {
		t.Fatalf("first message = %+v, want normalized user message", session.Messages[0])
	}
	if session.Messages[1].Role != "assistant" || session.Messages[1].Text != "The deploy notes are in session 42." {
		t.Fatalf("second message = %+v, want normalized assistant text block", session.Messages[1])
	}

	wantTime := time.Date(2026, 5, 11, 1, 2, 0, 0, time.UTC)
	if !session.Messages[0].Timestamp.Equal(wantTime) {
		t.Fatalf("timestamp = %s, want %s", session.Messages[0].Timestamp, wantTime)
	}
}

func TestParseClaudeTranscriptAcceptsArrayUserContent(t *testing.T) {
	input := strings.NewReader(`{"type":"user","sessionId":"s1","timestamp":"2026-05-11T01:02:00.000Z","message":{"role":"user","content":[{"type":"text","text":"hello"},{"type":"image","source":{}}]}}`)

	session, err := ParseClaudeTranscript(input, "s1.jsonl")
	if err != nil {
		t.Fatalf("ParseClaudeTranscript returned error: %v", err)
	}
	if len(session.Messages) != 1 || session.Messages[0].Text != "hello" {
		t.Fatalf("messages = %+v, want text-only user content", session.Messages)
	}
}

func TestParseClaudeTranscriptFallsBackToFilenameID(t *testing.T) {
	input := strings.NewReader(`{"type":"user","timestamp":"2026-05-11T01:02:00.000Z","message":{"role":"user","content":"hello"}}`)

	session, err := ParseClaudeTranscript(input, `C:\Users\me\.claude\projects\example-session.jsonl`)
	if err != nil {
		t.Fatalf("ParseClaudeTranscript returned error: %v", err)
	}
	if session.ID != "example-session" {
		t.Fatalf("session.ID = %q, want example-session", session.ID)
	}
}
