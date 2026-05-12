package source

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/LHYintheCode/agentdata/internal/model"
)

type hermesRecord struct {
	SessionID    string          `json:"session_id"`
	SessionIDAlt string          `json:"sessionId"`
	Project      string          `json:"project"`
	CWD          string          `json:"cwd"`
	Timestamp    string          `json:"timestamp"`
	Role         string          `json:"role"`
	Content      json.RawMessage `json:"content"`
	Message      hermesMessage   `json:"message"`
}

type hermesMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

func ParseHermesTranscript(r io.Reader, filename string) (model.Session, error) {
	session := model.Session{
		ID:     fallbackSessionID(filename),
		Source: "hermes",
	}

	scanner := newJSONLScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if lineNumber == 1 {
			line = strings.TrimPrefix(line, "\ufeff")
		}

		var record hermesRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return model.Session{}, fmt.Errorf("parse hermes line %d: %w", lineNumber, err)
		}
		if id := firstNonEmpty(record.SessionID, record.SessionIDAlt); id != "" {
			session.ID = id
		}
		if session.Project == "" {
			session.Project = firstNonEmpty(record.Project, record.CWD)
		}

		message, ok, err := hermesRecordMessage(record)
		if err != nil {
			return model.Session{}, fmt.Errorf("parse hermes line %d message: %w", lineNumber, err)
		}
		if ok {
			session.Messages = append(session.Messages, message)
		}
	}
	if err := scanner.Err(); err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func hermesRecordMessage(record hermesRecord) (model.Message, bool, error) {
	role := firstNonEmpty(record.Role, record.Message.Role)
	if role != "user" && role != "assistant" {
		return model.Message{}, false, nil
	}

	rawContent := record.Content
	if len(rawContent) == 0 {
		rawContent = record.Message.Content
	}
	text, err := claudeText(rawContent)
	if err != nil {
		return model.Message{}, false, err
	}
	if text == "" {
		return model.Message{}, false, nil
	}

	timestamp, err := parseTimestamp(record.Timestamp)
	if err != nil {
		return model.Message{}, false, err
	}
	return model.Message{
		Role:      role,
		Text:      text,
		Timestamp: timestamp,
	}, true, nil
}
