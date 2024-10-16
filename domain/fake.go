// // File: internal/app/app.go
package domain

//
//import (
//	"database/sql"
//	"github.com/bertoxic/graphql_chat/internal/app"
//	"log"
//	"yourproject/internal/auth"
//"yourproject/internal/handlers"
//"yourproject/internal/repository"
//"yourproject/internal/config"
//)
//
//type App struct {
//	config             *config.Config
//	authService        auth.AuthService
//	registrationHandler *handlers.RegistrationHandler
//	// Add other services and handlers as needed
//}
//
//func NewApp(cfg *config.Config) (*App, error) {
//	app := &App{
//		config: cfg,
//	}
//	if err := app.initialize(); err != nil {
//		return nil, err
//	}
//	return app, nil
//}
//
//func (a *App) initialize() error {
//	if err := a.initializeRepository(); err != nil {
//		return err
//	}
//	if err := a.initializeServices(); err != nil {
//		return err
//	}
//	if err := a.initializeHandlers(); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (a *App) initializeRepository() error {
//	// Initialize your database connection
//	db, err := initializeDatabase(a.config.DatabaseURL)
//	if err != nil {
//		return err
//	}
//
//	// Initialize your repositories
//	a.userRepo = repository.NewUserRepository(db)
//	return nil
//}
//
//func (a *App) initializeServices() error {
//	// Initialize your services
//	a.authService = auth.NewAuthService(a.userRepo)
//	return nil
//}
//
//func (a *App) initializeHandlers() error {
//	// Initialize your handlers
//	a.registrationHandler = handlers.NewRegistrationHandler(a.authService)
//	return nil
//}
//
//func (a *App) Run() error {
//	// Set up and start your HTTP server
//	// Use a.registrationHandler and other initialized components
//	// ...
//}
//
//// Utility function to initialize database
//func initializeDatabase(dbURL string) (*sql.Repo, error) {
//	// Initialize and return database connection
//}
//
//// File: cmd/server/main.go
//
//package main
//
//import (
//"log"
//"yourproject/internal/app"
//"yourproject/internal/config"
//)
//
//func main() {
//	cfg, err := config.Load()
//	if err != nil {
//		log.Fatalf("Failed to load config: %v", err)
//	}
//
//	application, err := app.NewApp(cfg)
//	if err != nil {
//		log.Fatalf("Failed to initialize application: %v", err)
//	}
//
//	if err := application.Run(); err != nil {
//		log.Fatalf("Error running application: %v", err)
//	}
//}
