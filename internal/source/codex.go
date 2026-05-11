package source

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/LHYintheCode/agentdata/internal/model"
)

type codexRecord struct {
	Timestamp string          `json:"timestamp"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
}

type codexSessionMeta struct {
	ID  string `json:"id"`
	CWD string `json:"cwd"`
}

type codexResponseItem struct {
	Role    string              `json:"role"`
	Content []codexContentBlock `json:"content"`
}

type codexContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func ParseCodexRollout(r io.Reader, filename string) (model.Session, error) {
	session := model.Session{
		ID:     fallbackSessionID(filename),
		Source: "codex",
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

		var record codexRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return model.Session{}, fmt.Errorf("parse codex line %d: %w", lineNumber, err)
		}

		switch record.Type {
		case "session_meta":
			if err := applyCodexSessionMeta(&session, record.Payload); err != nil {
				return model.Session{}, fmt.Errorf("parse codex line %d session_meta: %w", lineNumber, err)
			}
		case "response_item":
			message, ok, err := codexMessage(record)
			if err != nil {
				return model.Session{}, fmt.Errorf("parse codex line %d response_item: %w", lineNumber, err)
			}
			if ok {
				session.Messages = append(session.Messages, message)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func applyCodexSessionMeta(session *model.Session, payload json.RawMessage) error {
	var meta codexSessionMeta
	if err := json.Unmarshal(payload, &meta); err != nil {
		return err
	}
	if meta.ID != "" {
		session.ID = meta.ID
	}
	if meta.CWD != "" {
		session.Project = meta.CWD
	}
	return nil
}

func codexMessage(record codexRecord) (model.Message, bool, error) {
	var item codexResponseItem
	if err := json.Unmarshal(record.Payload, &item); err != nil {
		return model.Message{}, false, err
	}
	if item.Role != "user" && item.Role != "assistant" {
		return model.Message{}, false, nil
	}

	text := codexText(item.Content)
	if text == "" {
		return model.Message{}, false, nil
	}

	timestamp, err := parseTimestamp(record.Timestamp)
	if err != nil {
		return model.Message{}, false, err
	}
	return model.Message{
		Role:      item.Role,
		Text:      text,
		Timestamp: timestamp,
	}, true, nil
}

func codexText(blocks []codexContentBlock) string {
	parts := make([]string, 0, len(blocks))
	for _, block := range blocks {
		if strings.TrimSpace(block.Text) == "" {
			continue
		}
		parts = append(parts, block.Text)
	}
	return strings.Join(parts, "\n")
}

func fallbackSessionID(filename string) string {
	base := filepath.Base(filename)
	return strings.TrimSuffix(base, filepath.Ext(base))
}
