package auth

import (
	"context"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/bertoxic/graphqlChat/internal/models"

	"github.com/bertoxic/graphqlChat/internal/database"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	CreateUser(ctx context.Context, user RegistrationInput) (*models.UserDetails, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserDetails, error)
}

type UserRepo struct {
	DB database.DatabaseRepo
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (us *UserRepo) CreateUser(ctx context.Context, user RegistrationInput) (*models.UserDetails, error) {
	// us.DB.CreateUser()
	//TODO implement me
	userdetails, err := us.DB.CreateUser(ctx, &user)
	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", err)
	}

	return userdetails, nil
}

func (UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.UserDetails, error) {
	//TODO implement me
	return &models.UserDetails{}, nil
}
