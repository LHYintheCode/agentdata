package source

import (
	"strings"
	"testing"
	"time"
)

func TestParseOpenClawSessionsNormalizesMessages(t *testing.T) {
	input := strings.NewReader(`{
		"sessions": [
			{
				"id": "openclaw-session-1",
				"workspace": "D:\\go_project",
				"messages": [
					{"role": "user", "content": "Find deploy notes", "timestamp": "2026-05-12T01:02:00Z"},
					{"role": "assistant", "content": "Found them in the release session.", "timestamp": "2026-05-12T01:03:00Z"}
				]
			}
		]
	}`)

	sessions, err := ParseOpenClawSessions(input, "sessions.json")
	if err != nil {
		t.Fatalf("ParseOpenClawSessions returned error: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("len(sessions) = %d, want 1", len(sessions))
	}
	session := sessions[0]
	if session.ID != "openclaw-session-1" {
		t.Fatalf("session.ID = %q, want openclaw-session-1", session.ID)
	}
	if session.Source != "openclaw" {
		t.Fatalf("session.Source = %q, want openclaw", session.Source)
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

func TestParseOpenClawSessionsAcceptsTopLevelArray(t *testing.T) {
	input := strings.NewReader(`[{"id":"s1","messages":[{"role":"user","text":"hello"}]}]`)

	sessions, err := ParseOpenClawSessions(input, "sessions.json")
	if err != nil {
		t.Fatalf("ParseOpenClawSessions returned error: %v", err)
	}
	if len(sessions) != 1 || len(sessions[0].Messages) != 1 {
		t.Fatalf("sessions = %+v, want one session with one message", sessions)
	}
}
