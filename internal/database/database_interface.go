package database

import (
	"context"
	"github.com/bertoxic/graphqlChat/internal/models"
)

type DatabaseRepo interface {
	CreateUser(ctx context.Context, user models.RegistrationInput) (*models.UserDetails, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserDetails, error)
}
