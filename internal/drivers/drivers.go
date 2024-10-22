package drivers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
	dsn  string
}

func NewPostgresDB(ctx context.Context, dsn string) (*PostgresDB, error) {
	dbConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("can't parse postgress config ")
	}
	conn, err := pgxpool.ConnectConfig(ctx, dbConf)
	if err != nil {
		return nil, err
	}
	db := &PostgresDB{Pool: conn}
	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{Pool: conn, dsn: dsn}, nil

}

func (db *PostgresDB) Ping(ctx context.Context) error {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("can't ping the database: %v", err)
		return err
	}
	log.Print("postgres pinged and well connected🤗🥰")
	return nil
}
func (db *PostgresDB) Close() {
	db.Pool.Close()
}

func (db *PostgresDB) GetPoolConn() *pgxpool.Pool {
	return db.Pool
}

//func NewDatabase(ctx context.Context, dsn string, driver string) (Database, error) {
//	switch driver {
//	case "pgx":
//		apperr := utils.NewAppError(1000, "database nulls", errors.New("db failed not"))
//		fmt.Printf("%v", apperr.Details)
//
//		return NewPostgresDB(ctx, dsn)
//	case "mysql":
//		_, err := sql.Open("mysql", dsn)
//		if err != nil {
//			return nil, fmt.Errorf("failed to connect to mysql: %w", err)
//		}
//		// return &SqlDB{Repo: db}, nil
//		return nil, nil
//
//	default:
//		return nil, fmt.Errorf("unsupported driver: %s", driver)
//	}
//}

func (db *PostgresDB) Migrate() error {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("unable to determine current file location")
	}

	// Go one level up by using filepath.Dir() twice
	parentDir := filepath.Dir(filepath.Dir(currentFile))

	migrationPath := filepath.Join(parentDir, "migrations")

	// Use filepath.ToSlash to ensure forward slashes
	migrationURL := "file://" + filepath.ToSlash(migrationPath)

	// Initialize the migrationstance
	m, err := migrate.New(migrationURL, db.dsn)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No new migrations")
		} else {
			log.Printf("Migration failed: %v\n", err)
			return fmt.Errorf("error migrating: %w", err)
		}
	}

	//log.Println("migration done")
	return nil
}

// func (db *PostgresDB) Migrate() error {
// 	_, currentFile, _, ok := runtime.Caller(0)
// 	if !ok {
// 		return fmt.Errorf("unable to determine current file location")
// 	}

// 	// Go one level up by using filepath.Dir() twice
// 	parentDir := filepath.Dir(filepath.Dir(currentFile))

// 	migrationPath := filepath.Join(parentDir, "migrations")

// 	// Initialize the migration instance
// 	m, err := migrate.New(fmt.Sprintf("file:///%s", migrationPath), db.dsn)
// 	if err != nil {
// 		return fmt.Errorf("failed to initialize migrations: %w", err)
// 	}
// 	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
// 		return fmt.Errorf("error migrate up: %v", err)
// 	}
// 	log.Println("migration done")
// 	return nil
// }
