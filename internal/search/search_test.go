package search

import (
	"testing"
	"time"

	"agentdata/internal/model"
)

func TestMessagesFindsCaseInsensitiveMatches(t *testing.T) {
	sessions := []model.Session{
		{
			ID:      "s1",
			Source:  "codex",
			Project: "D:\\go_project",
			Messages: []model.Message{
				{Role: "user", Text: "Fix the deployment pipeline", Timestamp: time.Date(2026, 5, 11, 1, 0, 0, 0, time.UTC)},
				{Role: "assistant", Text: "Tests are passing", Timestamp: time.Date(2026, 5, 11, 1, 1, 0, 0, time.UTC)},
			},
		},
	}

	results := Messages(sessions, "DEPLOY")
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].SessionID != "s1" {
		t.Fatalf("SessionID = %q, want s1", results[0].SessionID)
	}
	if results[0].Message.Text != "Fix the deployment pipeline" {
		t.Fatalf("Message.Text = %q, want deployment message", results[0].Message.Text)
	}
}

func TestMessagesTrimsEmptyQuery(t *testing.T) {
	sessions := []model.Session{
		{ID: "s1", Messages: []model.Message{{Text: "anything"}}},
	}

	results := Messages(sessions, "   ")
	if len(results) != 0 {
		t.Fatalf("len(results) = %d, want 0", len(results))
	}
}
