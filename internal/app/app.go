package app

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"

	"github.com/bertoxic/graphqlChat/internal/auth"
	"github.com/bertoxic/graphqlChat/internal/config"
	"github.com/bertoxic/graphqlChat/internal/drivers"
)

type App struct {
	config      config.AppConfig
	DB          drivers.Database
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

	return app, nil
}

func (a *App) initialize() error {
	if a.config.DataBase == nil || a.config.DataBase.URL == "" {
		return fmt.Errorf("database configuration is missing or incomplete")
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
	database, err := drivers.NewDatabase(ctx, a.config.DataBase.URL, "pgx")
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	a.DB = database
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
