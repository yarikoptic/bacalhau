package model

import "time"

// we emit these to other nodes so they update their
// state locally and can emit events locally
type WebSocketEvent struct {
	// APIVersion of the Job
	APIVersion string `json:"APIVersion,omitempty"`

	JobEvent JobEvent `json:"JobEvent,omitempty"`

	EventTime time.Time `json:"EventTime,omitempty"`
}
