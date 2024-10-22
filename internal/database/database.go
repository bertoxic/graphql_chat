package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/drivers"
	"github.com/jackc/pgx/v4/pgxpool"
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
