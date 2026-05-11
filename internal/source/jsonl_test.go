package source

import (
	"strings"
	"testing"
	"time"
)

func TestParseJSONLSessionsNormalizesMessages(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		`{"source":"claude-code","project":"D:\\go_project","session_id":"s1","timestamp":"2026-05-11T01:02:03Z","role":"user","content":"Fix the build"}`,
		`{"source":"claude-code","project":"D:\\go_project","session_id":"s1","timestamp":"2026-05-11T01:03:04Z","role":"assistant","content":"The failing package is internal/source."}`,
	}, "\n"))

	sessions, err := ParseJSONLSessions(input)
	if err != nil {
		t.Fatalf("ParseJSONLSessions returned error: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("len(sessions) = %d, want 1", len(sessions))
	}

	session := sessions[0]
	if session.Source != "claude-code" {
		t.Fatalf("session.Source = %q, want claude-code", session.Source)
	}
	if session.Project != `D:\go_project` {
		t.Fatalf("session.Project = %q, want D:\\go_project", session.Project)
	}
	if session.ID != "s1" {
		t.Fatalf("session.ID = %q, want s1", session.ID)
	}
	if len(session.Messages) != 2 {
		t.Fatalf("len(session.Messages) = %d, want 2", len(session.Messages))
	}
	if session.Messages[0].Role != "user" || session.Messages[0].Text != "Fix the build" {
		t.Fatalf("first message = %+v, want normalized user content", session.Messages[0])
	}

	wantTime := time.Date(2026, 5, 11, 1, 2, 3, 0, time.UTC)
	if !session.Messages[0].Timestamp.Equal(wantTime) {
		t.Fatalf("timestamp = %s, want %s", session.Messages[0].Timestamp, wantTime)
	}
}

func TestParseJSONLSessionsRejectsInvalidJSON(t *testing.T) {
	_, err := ParseJSONLSessions(strings.NewReader(`{"source":"codex"`))
	if err == nil {
		t.Fatal("ParseJSONLSessions returned nil error for invalid JSON")
	}
}

func TestParseJSONLSessionsAcceptsUTF8BOM(t *testing.T) {
	input := strings.NewReader("\ufeff" + `{"source":"codex","project":"D:\\go_project","session_id":"s1","timestamp":"2026-05-11T01:02:03Z","role":"user","content":"Deploy the CLI"}`)

	sessions, err := ParseJSONLSessions(input)
	if err != nil {
		t.Fatalf("ParseJSONLSessions returned error: %v", err)
	}
	if len(sessions) != 1 || len(sessions[0].Messages) != 1 {
		t.Fatalf("sessions = %+v, want one session with one message", sessions)
	}
}
