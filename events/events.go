package events

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	Event struct {
		ID        primitive.ObjectID     `json:"_id" bson:"_id"`
		EventKind string                 `json:"event_kind" bson:"event_kind"`
		EventKey  string                 `json:"event_key" bson:"event_key"`
		MsgBody   map[string]interface{} `json:"msg_body" bson:"msg_body"`
		Processed bool                   `json:"processed" bson:"processed"`
	}
)
