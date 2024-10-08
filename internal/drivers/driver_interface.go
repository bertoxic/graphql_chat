package drivers

import "context"

type Database interface {
	Ping(ctx context.Context) error
	Close()
	Migrate() error
}
