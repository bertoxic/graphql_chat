package postgres

import (
	"context"
	"github.com/bertoxic/graphqlChat/internal/app"
	"github.com/bertoxic/graphqlChat/internal/drivers"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDBRepo struct {
	App *app.App
	DB  *pgxpool.Pool
}

func NewPostgresDBRepo(a *app.App, db *drivers.PostgresDB) *PostgresDBRepo {
	return &PostgresDBRepo{
		App: a,
		DB:  db.Pool,
	}
}
func (pr *PostgresDBRepo) CreateUser(ctx context.Context, user models.InputDetails) (*models.UserDetails, error) {
	return &models.UserDetails{}, nil
}
func (pr *PostgresDBRepo) GetUserByEmail(ctx context.Context, email string) (*models.UserDetails, error) {

	return &models.UserDetails{}, nil
}
