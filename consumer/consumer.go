package consumer

import (
	"context"
	"encoding/json"
	"fmt"
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
	Consumer struct {
		maxWorkerChans uint
		refreshTime    time.Duration
		handlers       map[string]Handler
	}

	Message struct {
		Key  string `json:"key"`
		Body string `json:"body"`
	}

	Handler func(ctx context.Context, msg Message) error
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

	workerChannels := make(chan Message, l.maxWorkerChans)
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
			var message Message

			if err := cursor.Decode(&message); err != nil {
				log.Fatal(err)
			}
			workerChannels <- message
			fmt.Printf("%+v\n", message)
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
func (l *Consumer) worker(ctx context.Context, msgChan <-chan Message) {
	for msg := range msgChan {
		_ = l.dispatcher(ctx, msg)
	}
}

func (l *Consumer) dispatcher(ctx context.Context, message Message) error {
	msg := new(Message)

	err := json.Unmarshal([]byte(message.Body), msg)
	if err != nil {
		return err
	}

	handlerFunc := l.handlers[message.Key]

	if handlerFunc == nil {
		zap.L().Info("handlerfunc is nil", zap.Error(err), zap.String("message", message.Key))
		return nil
	}

	if err := handlerFunc(ctx, message); err != nil {
		zap.L().Error("failed to handle message", zap.Error(err), zap.String("message", message.Key))
		return err
	}

	return nil
}
