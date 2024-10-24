package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/drivers"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

type Database interface {
	Ping(ctx context.Context) error
	Close()
	Migrate() error
	GetPoolConn() *pgxpool.Pool
}

func NewDatabase(ctx context.Context, dsn string, driver string) (Database, error) {
	switch driver {
	case "pgx":
		return drivers.NewPostgresDB(ctx, dsn)
	case "redis":
		return drivers.NewRedisDB(ctx, dsn)
	case "mysql":
		_, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to mysql: %w", err)
		}
		// return &SqlDB{Repo: db}, nil
		return nil, nil

	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}
}

type RedisClient struct {
	Client *redis.Client
}
type RedisRepo interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) string
}

func (r RedisClient) Set(ctx context.Context, key string, value string, duration time.Duration) error {
	err := r.Client.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		fmt.Println("Error setting key:", err)
		return err
	}
	return nil
}

func (r RedisClient) Get(ctx context.Context, key string) (string, error) {
	res, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("Error getting value:", err)
		return "", err
	}
	return res, nil
}

func NewRedisClient(ctx context.Context, client *redis.Client) (*RedisClient, error) {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return nil, err
	}
	fmt.Println("Redis connected:", pong)
	return &RedisClient{Client: client}, nil
}
