package source

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/LHYintheCode/agentdata/internal/model"
)

type jsonlRecord struct {
	Source    string `json:"source"`
	Project   string `json:"project"`
	SessionID string `json:"session_id"`
	Timestamp string `json:"timestamp"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	Text      string `json:"text"`
}

func ParseJSONLSessions(r io.Reader) ([]model.Session, error) {
	scanner := bufio.NewScanner(r)
	sessionsByID := make(map[string]*model.Session)
	order := make([]string, 0)
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

		var record jsonlRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, fmt.Errorf("parse jsonl line %d: %w", lineNumber, err)
		}

		sessionID := record.SessionID
		if sessionID == "" {
			sessionID = "default"
		}

		session, ok := sessionsByID[sessionID]
		if !ok {
			session = &model.Session{
				ID:      sessionID,
				Source:  record.Source,
				Project: record.Project,
			}
			sessionsByID[sessionID] = session
			order = append(order, sessionID)
		}

		timestamp, err := parseTimestamp(record.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("parse jsonl line %d timestamp: %w", lineNumber, err)
		}

		text := record.Content
		if text == "" {
			text = record.Text
		}

		session.Messages = append(session.Messages, model.Message{
			Role:      record.Role,
			Text:      text,
			Timestamp: timestamp,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sessions := make([]model.Session, 0, len(order))
	for _, id := range order {
		sessions = append(sessions, *sessionsByID[id])
	}
	return sessions, nil
}

func parseTimestamp(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339Nano, value)
}
