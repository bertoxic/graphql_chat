package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/database/postgres"
	"github.com/bertoxic/graphqlChat/internal/drivers"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/bertoxic/graphqlChat/internal/handlers"
	"github.com/bertoxic/graphqlChat/internal/jwt"
	"github.com/bertoxic/graphqlChat/internal/render"
	"github.com/bertoxic/graphqlChat/internal/user"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/bertoxic/graphqlChat/internal/auth"
	"github.com/bertoxic/graphqlChat/internal/database"
)

type App struct {
	Config   *config.AppConfig
	DB       database.DatabaseRepo
	Services *ServicesContainer
}
type ServicesContainer struct {
	AuthService     auth.AuthService
	UserAuthService auth.UserRepository
	UserService     *user.Service
}

//
//func getLongPathName(shortPath string) (string, error) {
//	longPath := make([]uint16, windows.MAX_LONG_PATH)
//	n, err := windows.GetLongPathName(&windows.StringToUTF16(shortPath)[0], &longPath[0], windows.MAX_LONG_PATH)
//	if err != nil {
//		return "", err
//	}
//	return windows.UTF16ToString(longPath[:n]), nil
//}

//nolint:funlen
func NewApp(ctx context.Context, cfg *config.AppConfig) (*App, error) {

	app := &App{
		Config:   cfg,
		Services: &ServicesContainer{},
	}
	app.Config.InProduction = true
	if err := app.initialize(); err != nil {
		return nil, err
	}

	if err := app.initializeDB(ctx); err != nil {
		return nil, err
	}
	if err := app.initializeServices(); err != nil {
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
	if a.Config.DataBaseINFO == nil || a.Config.DataBaseINFO.URL == "" {
		return errorx.New(errorx.ErrCodeValidation, "Data configuration is invalid or missing", fmt.Errorf("could not load database Config"))
	}
	return nil
}

func (a *App) initializeServices() error {
	//initialize all my Services here
	authUserRepo := auth.NewUserRepo(a.DB)
	tokenService := jwt.NewTokenService(a.Config)
	a.Services.AuthService = auth.NewAuthService(authUserRepo, tokenService)
	a.Services.UserAuthService = auth.NewUserRepo(a.DB)
	userRepo := user.NewUserRepo(a.DB)
	userService := user.NewService(userRepo)
	a.Services.UserService = userService
	return nil
}

func (a *App) initializeDB(ctx context.Context) error {
	newDatabase, err := database.NewDatabase(ctx, a.Config.DataBaseINFO.URL, "pgx")
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
	a.DB = postgres.NewPostgresDBRepo(a.Config, db)
	return nil
}

func (a *App) initializeHandlers() error {
	dbrepo := handlers.NewRepository(a.Config, a.DB)
	handlers.NewRepo(dbrepo)
	return nil
}

func (a *App) initializeRender() error {
	render.NewRenderer(a.Config)
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		var appErr *errorx.AppError
		ok := errors.As(err, &appErr)
		if !ok {
			return err
		}
		fmt.Printf("%s", appErr.Details)
		return err
	}
	a.Config.TemplateCache = templateCache
	return nil
}

//func LoadEnv() error {
//	// Check if environment variables are already set (e.g., in production)
//	// if os.Getenv("DB_HOST") != "" {
//	//	fmt.Println("Environment variables already set, skipping .env file loading")
//	//	return nil
//	//}
//
//	// List of possible locations for the .env file
//	possiblePaths := []string{
//		".env",
//		"../.env",
//		"../../.env",
//		"../../../.env",
//	}
//
//	// Get the executable path
//	ex, err := os.Executable()
//	if err == nil {
//		exePath := filepath.Dir(ex)
//		possiblePaths = append(possiblePaths,
//			filepath.Join(exePath, ".env"),
//			filepath.Join(exePath, "../.env"),
//		)
//	}
//
//	// Try to load .env from the possible locations
//	for _, path := range possiblePaths {
//		_, _ = filepath.Abs(path)
//		// fmt.Printf("Trying to load .env from: %s\n", absPath)
//		err := godotenv.Load(path)
//		if err == nil {
//			//fmt.Printf("Successfully loaded .env from: %s\n", absPath)
//			return nil
//		}
//	}
//
//	return fmt.Errorf("unable to load .env file from any location")
//}

func LoadEnvInProd() error {
	// Check if environment variables are already set (e.g., in production)
	if os.Getenv("DB_HOST") != "" {
		fmt.Println("Environment variables already set, skipping .env file loading")
		return nil
	}

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
		absPath, _ := filepath.Abs(path)
		// fmt.Printf("Trying to load .env from: %s\n", absPath)
		err := godotenv.Load(absPath)
		if err == nil {
			// fmt.Printf("Successfully loaded .env from: %s\n", absPath)
			return nil
		}
	}

	return fmt.Errorf("unable to load .env file from any location")
}
