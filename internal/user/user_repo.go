package user

import (
	"context"
	"github.com/bertoxic/graphqlChat/internal/app"
	"github.com/bertoxic/graphqlChat/internal/drivers"
)

type UserRepo interface {
}

type Repository struct {
	DB  drivers.Database
	app *app.App
}

func NewUserRepo(db drivers.Database, a *app.App) *Repository {
	return &Repository{
		DB:  db,
		app: a,
	}
}

func (us *Repository) GetUserByEmail(ctx context.Context, email string) {

}
