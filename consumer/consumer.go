package consumer

import (
	"context"
	"errors"
	"fmt"
	"github.com/tejiriaustin/narx_api/events"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/tejiriaustin/narx_api/repository"
)

const (
	defaultMaxWorkers uint = 10
)

type (
	Fetcher interface {
		Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	}
	Updater interface {
		UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	}

	Consumer struct {
		maxWorkerChans uint
		refreshTime    time.Duration
		handlers       map[string]Handler
		updater        Updater
	}

	Handler func(ctx context.Context, msg events.Event) error

	Options func(*Consumer)
)

func newConsumer() *Consumer {
	return &Consumer{
		maxWorkerChans: defaultMaxWorkers,
		refreshTime:    15 * time.Second,
		handlers:       make(map[string]Handler),
	}
}

func NewConsumer(opts ...Options) *Consumer {
	l := newConsumer()

	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *Consumer) SetHandler(key string, handler Handler) *Consumer {
	l.handlers[strings.ToUpper(key)] = handler
	return l
}

func (l *Consumer) ListenAndServe(ctx context.Context, pubSub Fetcher) {
	log.Print("initializing  Consumer...")

	workerChannels := make(chan events.Event, l.maxWorkerChans)
	defer func() {
		close(workerChannels)
	}()

	log.Print("starting workers\n", "maxWorkers: ", l.maxWorkerChans)

	for i := uint(0); i < l.maxWorkerChans; i++ {
		go l.worker(ctx, workerChannels)
	}

	for {
		zap.L().Info("pulling messages...")
		cursor, err := pubSub.Find(ctx, repository.NewQueryFilter().AddFilter("processed", false))
		if err != nil {
			zap.L().Error("failed to receive message", zap.Error(err))
		}

		for cursor.Next(ctx) {
			var message events.Event

			if err = cursor.Decode(&message); err != nil {
				log.Fatal(err)
			}
			workerChannels <- message
		}
		if err = cursor.Err(); err != nil {
			log.Fatal(err)
		}

		time.Sleep(l.refreshTime)
	}
}

// worker listens on the msgChan for incoming messages.
// as messages become available over the channel, it looks over the map of configured handlers and routes messages by the key.
// If no handlers are configured and a default handler has been set, the message is sent there.
// else it logs and continues.
func (l *Consumer) worker(ctx context.Context, msgChan <-chan events.Event) {
	for msg := range msgChan {
		err := l.dispatcher(ctx, msg)
		if err != nil {
			zap.L().Error(err.Error(), zap.String("message", msg.EventKind))
		}
	}
}

func (l *Consumer) dispatcher(ctx context.Context, message events.Event) error {
	handlerFunc := l.handlers[message.EventKey]

	if handlerFunc == nil {
		zap.L().Error("handler func is nil", zap.String("message", message.EventKind))
		return errors.New("handler func is nil")
	}
	fmt.Println("handler func")

	if err := handlerFunc(ctx, message); err != nil {
		zap.L().Error("failed to handle message", zap.Error(err), zap.String("message", message.EventKind))
		return err
	}

	updates := map[string]interface{}{
		"$sets": map[string]interface{}{
			"processed": true,
		},
	}

	_, err := l.updater.UpdateOne(ctx, repository.NewQueryFilter().AddFilter("_id", message.ID), updates)
	if err != nil {
		zap.L().Error("failed to update message processed", zap.Error(err), zap.String("message", message.EventKind))
		return err
	}
	return nil
}
