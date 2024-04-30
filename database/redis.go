package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	SubscriptionChannel = "SUBSCRIPTION_CHANNEL"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(dsn, password, port string) (*RedisClient, error) {
	log.Println(" connecting to redis database...")

	redisDbClient := redis.NewClient(&redis.Options{
		Password: password,
		Addr:     dsn,
		DB:       0,
	})

	if err := redisDbClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Error connecting to redis: %v", err)
		return nil, err
	}

	return &RedisClient{
		rdb: redisDbClient,
	}, nil
}

func (rdb *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	return rdb.rdb.Set(ctx, key, value, 0).Err()
}

func (rdb *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := rdb.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (rdb *RedisClient) Disconnect(ctx context.Context) error {
	err := rdb.rdb.Close()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *RedisClient) Subscribe(ctx context.Context) *redis.PubSub {
	return rdb.rdb.Subscribe(ctx, SubscriptionChannel)
}

func (rdb *RedisClient) Publish(ctx context.Context, msg interface{}) error {
	return rdb.rdb.Publish(ctx, SubscriptionChannel, msg).Err()
}
