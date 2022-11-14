package model

import "time"

// we emit these to other nodes so they update their
// state locally and can emit events locally
type EventEnvelope struct {
	// APIVersion of the Job
	APIVersion string `json:"APIVersion,omitempty"`

	JobEvent JobEvent `json:"JobEvent,omitempty"`

	EventTime time.Time `json:"EventTime,omitempty"`
}

// Create a new EventEnvelope with current API version and time
func NewEventEnvelope(je JobEvent) EventEnvelope {
	return EventEnvelope{
		APIVersion: je.APIVersion,
		JobEvent:   je,
		EventTime:  time.Now(),
	}
}
