package slack

import (
	"encoding/json"
)

// EventsAPIEvent represents a Slack Events API event
type EventsAPIEvent struct {
	Type       string          `json:"type"`
	Event      json.RawMessage `json:"event"`
	EventID    string          `json:"event_id"`
	TeamID     string          `json:"team_id"`
	APIAppID   string          `json:"api_app_id"`
	Data       interface{}     `json:"-"`
	EventTime  int             `json:"event_time"`
	InnerEvent struct {
		Data    interface{} `json:"-"`
		Type    string      `json:"type"`
		EventID string      `json:"-"`
	} `json:"-"`
}