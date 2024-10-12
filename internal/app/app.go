package app

import (
	"context"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/database/postgres"
	"github.com/bertoxic/graphqlChat/internal/drivers"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/bertoxic/graphqlChat/internal/handlers"
	"github.com/bertoxic/graphqlChat/internal/render"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"golang.org/x/sys/windows"

	"github.com/bertoxic/graphqlChat/internal/auth"
	"github.com/bertoxic/graphqlChat/internal/database"
	"github.com/bertoxic/graphqlChat/pkg/config"
)

type App struct {
	config      *config.AppConfig
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
func NewApp(ctx context.Context, cfg *config.AppConfig) (*App, error) {

	app := &App{
		config: cfg,
	}
	if err := app.initialize(); err != nil {
		return nil, err
	}

	if err := app.initializeDB(ctx); err != nil {
		return nil, err
	}
	if err := app.initializeRender(); err != nil {
		return nil, err
	}
	if err := app.initializeHandlers(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initialize() error {
	if a.config.DataBaseINFO == nil || a.config.DataBaseINFO.URL == "" {
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
	newDatabase, err := database.NewDatabase(ctx, a.config.DataBaseINFO.URL, "pgx")
	if err != nil {
		return fmt.Errorf("failed to initialize newDatabase: %w", err)
	}
	err = newDatabase.Migrate()
	if err != nil {
		log.Fatalf("unable to run newDatabase migration: %v", err)
	}
	db, ok := newDatabase.(*drivers.PostgresDB)
	if !ok {
		return errorx.New(errorx.ErrCodeInternal, "the type assertion for databases failed", errorx.ErrDatabase)
	}
	a.DB = postgres.NewPostgresDBRepo(a.config, db)
	return nil
}

func (a *App) initializeHandlers() error {
	dbrepo := handlers.NewRepository(a.config, a.DB)
	handlers.NewRepo(dbrepo)
	return nil
}
func (a *App) initializeRender() error {
	render.NewRenderer(a.config)
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		appErr, ok := err.(*errorx.AppError)
		if !ok {
			return err
		}
		fmt.Printf("%s", appErr.Details)
		return err
	}
	a.config.TemplateCache = templateCache
	return nil
}

func LoadEnv() error {
	// Check if environment variables are already set (e.g., in production)
	// if os.Getenv("DB_HOST") != "" {
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
		// fmt.Printf("Trying to load .env from: %s\n", absPath)
		err := godotenv.Load(path)
		if err == nil {
			//fmt.Printf("Successfully loaded .env from: %s\n", absPath)
			return nil
		}
	}

	return fmt.Errorf("unable to load .env file from any location")
}
