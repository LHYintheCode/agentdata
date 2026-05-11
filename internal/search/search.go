package search

import (
	"strings"

	"github.com/LHYintheCode/agentdata/internal/model"
)

type Result struct {
	SessionID string
	Source    string
	Project   string
	Message   model.Message
}

func Messages(sessions []model.Session, query string) []Result {
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	if normalizedQuery == "" {
		return nil
	}

	results := make([]Result, 0)
	for _, session := range sessions {
		for _, message := range session.Messages {
			if strings.Contains(strings.ToLower(message.Text), normalizedQuery) {
				results = append(results, Result{
					SessionID: session.ID,
					Source:    session.Source,
					Project:   session.Project,
					Message:   message,
				})
			}
		}
	}
	return results
}
