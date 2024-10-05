package app

import (
	"github.com/bertoxic/graphqlChat/internal/auth"
	"github.com/bertoxic/graphqlChat/internal/config"
)

type App struct {
	config      config.AppConfig
	authService auth.AuthService
}

func NewApp(cfg config.AppConfig) (*App, error) {
	app := &App{
		config: cfg,
	}
	if err := app.initialize(); err != nil {
		return nil, err
	}
	return app, nil
}

func (a App) initialize() error {

	return nil
}

func (a *App) initializeServices() {
	//initialize all my services here
	userRepo := auth.NewUserRepo()
	a.authService = auth.NewAuthService(userRepo)

}
