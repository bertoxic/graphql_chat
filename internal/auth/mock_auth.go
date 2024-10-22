package auth

import (
	"context"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user RegistrationInput) (models.UserDetails, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.UserDetails), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (models.UserDetails, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.UserDetails), args.Error(1)
}
