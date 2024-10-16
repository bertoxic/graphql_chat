package auth

import (
	"context"
	"fmt"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/bertoxic/graphqlChat/internal/database"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.RegistrationInput) (*models.UserDetails, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserDetails, error)
}

type UserRepo struct {
	DB database.DatabaseRepo
}

func NewUserRepo(db database.DatabaseRepo) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (us *UserRepo) CreateUser(ctx context.Context, user *models.RegistrationInput) (*models.UserDetails, error) {

	// us.Repo.CreateUser()
	//TODO implement me
	userDetails, err := us.DB.CreateUser(ctx, *user)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", err)
	}

	return userDetails, nil
}
func createUser(ctx context.Context, pgx *pgxpool.Pool) error {
	return nil
}
func (us *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.UserDetails, error) {
	userDetails, err := us.DB.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", err)
		fmt.Printf("%s", err)
	}

	return userDetails, nil
}
