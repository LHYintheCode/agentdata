package exporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/LHYintheCode/agentdata/internal/model"
)

type JSONLRecord struct {
	Source    string `json:"source"`
	Project   string `json:"project"`
	SessionID string `json:"session_id"`
	Timestamp string `json:"timestamp"`
	Role      string `json:"role"`
	Text      string `json:"text"`
}

func JSONL(w io.Writer, sessions []model.Session) error {
	encoder := json.NewEncoder(w)
	for _, session := range sessions {
		for _, message := range session.Messages {
			record := JSONLRecord{
				Source:    session.Source,
				Project:   session.Project,
				SessionID: session.ID,
				Timestamp: message.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
				Role:      message.Role,
				Text:      message.Text,
			}
			if err := encoder.Encode(record); err != nil {
				return err
			}
		}
	}
	return nil
}

func Markdown(w io.Writer, sessions []model.Session) error {
	for _, session := range sessions {
		if _, err := fmt.Fprintf(w, "# %s\n\nSource: %s\n\nProject: %s\n\n", session.ID, session.Source, session.Project); err != nil {
			return err
		}
		for _, message := range session.Messages {
			if _, err := fmt.Fprintf(w, "## %s\n\n%s\n\n", message.Role, message.Text); err != nil {
				return err
			}
		}
	}
	return nil
}
