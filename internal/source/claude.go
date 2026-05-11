package source

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/LHYintheCode/agentdata/internal/model"
)

type claudeRecord struct {
	Type      string        `json:"type"`
	SessionID string        `json:"sessionId"`
	CWD       string        `json:"cwd"`
	Timestamp string        `json:"timestamp"`
	Message   claudeMessage `json:"message"`
}

type claudeMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type claudeContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func ParseClaudeTranscript(r io.Reader, filename string) (model.Session, error) {
	session := model.Session{
		ID:     fallbackSessionID(filename),
		Source: "claude",
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

		var record claudeRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return model.Session{}, fmt.Errorf("parse claude line %d: %w", lineNumber, err)
		}
		if record.Type != "user" && record.Type != "assistant" {
			continue
		}
		if record.SessionID != "" {
			session.ID = record.SessionID
		}
		if session.Project == "" && record.CWD != "" {
			session.Project = record.CWD
		}

		message, ok, err := claudeRecordMessage(record)
		if err != nil {
			return model.Session{}, fmt.Errorf("parse claude line %d message: %w", lineNumber, err)
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

func claudeRecordMessage(record claudeRecord) (model.Message, bool, error) {
	role := record.Message.Role
	if role == "" {
		role = record.Type
	}
	if role != "user" && role != "assistant" {
		return model.Message{}, false, nil
	}

	text, err := claudeText(record.Message.Content)
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

func claudeText(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return "", nil
	}

	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text, nil
	}

	var blocks []claudeContentBlock
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return "", err
	}

	parts := make([]string, 0, len(blocks))
	for _, block := range blocks {
		if block.Type != "" && block.Type != "text" {
			continue
		}
		if strings.TrimSpace(block.Text) == "" {
			continue
		}
		parts = append(parts, block.Text)
	}
	return strings.Join(parts, "\n"), nil
}

func IsClaudeTranscriptPath(path string) bool {
	if !strings.EqualFold(filepath.Ext(path), ".jsonl") {
		return false
	}
	return !strings.Contains(strings.ToLower(filepath.ToSlash(path)), "/subagents/")
}
