package events

type (
	Event struct {
		EventKind string                 `json:"event_kind"`
		EventKey  string                 `json:"event_key"`
		MsgBody   map[string]interface{} `json:"msg_body"`
	}
)
