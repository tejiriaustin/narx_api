package publisher

import (
	"context"
	"github.com/tejiriaustin/narx_api/events"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type (
	Inserter interface {
		InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	}

	Publisher struct {
		m        sync.Mutex
		inserter Inserter
	}

	PublishInterface interface {
		Publish(ctx context.Context, key, kind string, message map[string]interface{}) error
	}
)

func newPublisher(inserter Inserter) *Publisher {
	return &Publisher{
		m:        sync.Mutex{},
		inserter: inserter,
	}
}
func NewPublisher(inserter Inserter) PublishInterface {
	return newPublisher(inserter)
}

func (p *Publisher) Publish(ctx context.Context, key, kind string, message map[string]interface{}) error {
	id := message["id"].(string)
	event := events.Event{
		ID:        id,
		EventKind: kind,
		EventKey:  key,
		MsgBody:   message,
		Processed: false,
	}

	_, err := p.inserter.InsertOne(ctx, event)
	return err
}
