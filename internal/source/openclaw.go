package source

import (
	"encoding/json"
	"io"
	"path/filepath"

	"github.com/LHYintheCode/agentdata/internal/model"
)

type openClawEnvelope struct {
	Sessions []openClawSession `json:"sessions"`
}

type openClawSession struct {
	ID        string            `json:"id"`
	SessionID string            `json:"session_id"`
	Key       string            `json:"key"`
	Workspace string            `json:"workspace"`
	Project   string            `json:"project"`
	CWD       string            `json:"cwd"`
	Messages  []openClawMessage `json:"messages"`
}

type openClawMessage struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

func ParseOpenClawSessions(r io.Reader, filename string) ([]model.Session, error) {
	decoder := json.NewDecoder(r)
	var raw json.RawMessage
	if err := decoder.Decode(&raw); err != nil {
		return nil, err
	}

	var envelope openClawEnvelope
	if err := json.Unmarshal(raw, &envelope); err == nil && len(envelope.Sessions) > 0 {
		return normalizeOpenClawSessions(envelope.Sessions, filename)
	}

	var sessions []openClawSession
	if err := json.Unmarshal(raw, &sessions); err != nil {
		return nil, err
	}
	return normalizeOpenClawSessions(sessions, filename)
}

func normalizeOpenClawSessions(items []openClawSession, filename string) ([]model.Session, error) {
	sessions := make([]model.Session, 0, len(items))
	for _, item := range items {
		session := model.Session{
			ID:      firstNonEmpty(item.ID, item.SessionID, item.Key, fallbackSessionID(filename)),
			Source:  "openclaw",
			Project: firstNonEmpty(item.Workspace, item.Project, item.CWD),
		}

		for _, message := range item.Messages {
			if message.Role != "user" && message.Role != "assistant" {
				continue
			}
			text := firstNonEmpty(message.Content, message.Text)
			if text == "" {
				continue
			}
			timestamp, err := parseTimestamp(message.Timestamp)
			if err != nil {
				return nil, err
			}
			session.Messages = append(session.Messages, model.Message{
				Role:      message.Role,
				Text:      text,
				Timestamp: timestamp,
			})
		}
		if len(session.Messages) > 0 {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func IsOpenClawSessionsPath(path string) bool {
	return filepath.Base(path) == "sessions.json"
}
