package publisher

import (
	"context"
	"encoding/json"
	"github.com/tejiriaustin/narx_api/events"
	"sync"

	"github.com/tejiriaustin/narx_api/database"
)

type (
	Publisher struct {
		m      sync.Mutex
		client *database.RedisClient
	}

	PublishInterface interface {
		Publish(ctx context.Context, key string, message map[string]interface{}) error
	}
)

func newPublisher(client *database.RedisClient) *Publisher {
	return &Publisher{
		m:      sync.Mutex{},
		client: client,
	}
}
func NewPublisher(client *database.RedisClient) PublishInterface {
	return newPublisher(client)
}

func (p *Publisher) Publish(ctx context.Context, key string, message map[string]interface{}) error {
	event := events.Event{
		EventKey: key,
		MsgBody:  message,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, string(eventBytes))
}
