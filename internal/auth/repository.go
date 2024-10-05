package auth

import (
	"context"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	CreateUser(ctx context.Context, user RegistrationInput) (UserDetails, error)
	GetUserByEmail(ctx context.Context, email string) (UserDetails, error)
}

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (UserRepo) CreateUser(ctx context.Context, user RegistrationInput) (UserDetails, error) {
	//TODO implement me
	return UserDetails{}, nil
}

func (UserRepo) GetUserByEmail(ctx context.Context, email string) (UserDetails, error) {
	//TODO implement me
	return UserDetails{}, nil
}
