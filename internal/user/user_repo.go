package user

import (
	"context"
	"github.com/bertoxic/graphqlChat/internal/app"
	"github.com/bertoxic/graphqlChat/internal/database"
)

type UserRepo interface {
}

type Repository struct {
	DB  database.DatabaseRepo
	app *app.App
}

func NewUserRepo(db database.DatabaseRepo, a *app.App) *Repository {
	return &Repository{
		DB:  db,
		app: a,
	}
}

func (us *Repository) GetUserByEmail(ctx context.Context, email string) {

}
