package app

import (
	"context"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/database"
	"github.com/bertoxic/graphqlChat/internal/database/postgres"
	"github.com/bertoxic/graphqlChat/internal/drivers"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/bertoxic/graphqlChat/internal/handlers"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"github.com/joho/godotenv"
	"golang.org/x/sys/windows"
	"log"
	"os"
	"path/filepath"

	"github.com/bertoxic/graphqlChat/internal/auth"
)

type App struct {
	config      config.AppConfig
	DB          database.DatabaseRepo
	authService auth.AuthService
}

func getLongPathName(shortPath string) (string, error) {
	longPath := make([]uint16, windows.MAX_LONG_PATH)
	n, err := windows.GetLongPathName(&windows.StringToUTF16(shortPath)[0], &longPath[0], windows.MAX_LONG_PATH)
	if err != nil {
		return "", err
	}
	return windows.UTF16ToString(longPath[:n]), nil
}

//nolint:funlen
func NewApp(ctx context.Context, cfg config.AppConfig) (*App, error) {

	app := &App{
		config: cfg,
	}
	if err := app.initialize(); err != nil {
		return nil, err
	}
	if err := app.initializeDB(ctx); err != nil {
		return nil, err
	}
	if err := app.initializeHandlers(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initialize() error {
	if a.config.DataBase == nil || a.config.DataBase.URL == "" {
		return errorx.New(errorx.ErrCodeValidation, "Data configuration is invalid or missing", fmt.Errorf("could not load database config"))
	}
	return nil
}

func (a *App) initializeServices() error {
	//initialize all my services here
	userRepo := auth.NewUserRepo()
	a.authService = auth.NewAuthService(userRepo)

	return nil
}

func (a *App) initializeDB(ctx context.Context) error {
	database, err := database.NewDatabase(ctx, a.config.DataBase.URL, "pgx")
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	err = database.Migrate()
	if err != nil {
		log.Fatalf("unable to run database migration: %v", err)
	}
	db, ok := database.(*drivers.PostgresDB)
	if ok {
		a.DB = postgres.NewPostgresDBRepo(a, db)
	}

	return nil
}

func (a *App) initializeHandlers(ctx context.Context, app *App) error {
	dbrepo := handlers.NewRepository(app, a.DB)
	handlers.NewRepo(dbrepo)
	return nil
}
func LoadEnv() error {
	// Check if environment variables are already set (e.g., in production)
	//if os.Getenv("DB_HOST") != "" {
	//	fmt.Println("Environment variables already set, skipping .env file loading")
	//	return nil
	//}

	// List of possible locations for the .env file
	possiblePaths := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
	}

	// Get the executable path
	ex, err := os.Executable()
	if err == nil {
		exePath := filepath.Dir(ex)
		possiblePaths = append(possiblePaths,
			filepath.Join(exePath, ".env"),
			filepath.Join(exePath, "../.env"),
		)
	}

	// Try to load .env from the possible locations
	for _, path := range possiblePaths {
		_, _ = filepath.Abs(path)
		//fmt.Printf("Trying to load .env from: %s\n", absPath)
		err := godotenv.Load(path)
		if err == nil {
			//fmt.Printf("Successfully loaded .env from: %s\n", absPath)
			return nil
		}
	}

	return fmt.Errorf("unable to load .env file from any location")
}
