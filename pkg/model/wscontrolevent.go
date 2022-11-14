package model

import "fmt"

//go:generate stringer -type=WSControlEventType --trimprefix=WSControlEvent
type WSControlEventType int

// type wsControlEvent struct {
// 	ActionType string `json:"action_type"` // subscribe, unsubscribe
// 	TaskID     string `json:"task_id"`
// }

const (
	wsControlEventUnknown WSControlEventType = iota // must be first

	// WebService subscribe to a feed
	WSControlEventSubscribe

	// WebService unsubscribe to a feed
	WSControlEventUnsubscribe

	wsControlEventDone // must be last
)

func ParseWSControlEventType(str string) (WSControlEventType, error) {
	for typ := wsControlEventUnknown + 1; typ < wsControlEventDone; typ++ {
		if equal(typ.String(), str) {
			return typ, nil
		}
	}

	return wsControlEventUnknown, fmt.Errorf(
		"executor: unknown web service control event type '%s'", str)
}

func WSControlEventTypes() []WSControlEventType {
	var res []WSControlEventType
	for typ := wsControlEventUnknown + 1; typ < wsControlEventDone; typ++ {
		res = append(res, typ)
	}

	return res
}

func (wsce WSControlEventType) MarshalText() ([]byte, error) {
	return []byte(wsce.String()), nil
}

func (wsce *WSControlEventType) UnmarshalText(text []byte) (err error) {
	name := string(text)
	*wsce, err = ParseWSControlEventType(name)
	return
}
