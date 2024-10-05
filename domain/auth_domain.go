package domain

//
//import (
//	"context"
//	"fmt"
//	"github.com/bertoxic/graphql_chat/internal/auth"
//	"github.com/bertoxic/graphql_chat/packages"
//	"github.com/bertoxic/graphql_chat/packages/utils"
//	"golang.org/x/crypto/bcrypt"
//)
//
//// Domain-specific types
//type RegistrationInput struct {
//	Email    string
//	Username string
//	Password string
//}
//
//type AuthResponse struct {
//	AccessToken string
//}
//
//// UserRepository defines the interface for user-related database operations
//type UserRepository interface {
//	CreateUser(ctx context.Context, user packages.User) error
//	GetUserByEmail(ctx context.Context, email, username string) (bool, error)
//}
//
//// User represents a user in the domain layer
//type User struct {
//	ID       string
//	Email    string
//	Username string
//	Password string
//}
//
//// AuthService defines the authentication service for the domain layer
//type AuthService struct {
//	userRepo UserRepository
//}
//
//// NewAuthService creates a new AuthService
//func NewAuthService(ur UserRepository) *AuthService {
//	return &AuthService{
//		userRepo: ur,
//	}
//}
//
//// Register handles user registration
//func (as *AuthService) Register(ctx context.Context, input auth.RegistrationInput) (*AuthResponse, error) {
//	// Validate input
//	if err := validateRegistrationInput(input); err != nil {
//		return nil, utils.ErrInvalidInput
//	}
//
//	// Check if user already exists
//	exists, err := as.userRepo.GetUserByEmail(ctx, input.Email, input.Username)
//	if err != nil {
//		return nil, err
//	}
//	if exists {
//		return nil, utils.ErrUserExists
//	}
//
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
//	if err != nil {
//		return nil, fmt.Errorf("%w: failed to hash password", utils.ErrInternalServer)
//	}
//
//	// TODO: Implement user storage and token generation
//	_ = hashedPassword // Placeholder to use hashedPassword
//
//	user := packages.User{
//		Email:    input.Email,
//		Username: input.Username,
//		Password: string(hashedPassword),
//	}
//	if err := as.userRepo.CreateUser(ctx, user); err != nil {
//		return nil, err
//	}
//
//	// Generate access token (implementation details omitted)
//	accessToken, err := generateAccessToken(user)
//	if err != nil {
//		return nil, err
//	}
//
//	return &AuthResponse{
//		AccessToken: accessToken,
//	}, nil
//}
//
//func validateRegistrationInput(input auth.RegistrationInput) error {
//
//	return nil
//}
//
//func generateAccessToken(user packages.User) (string, error) {
//	return "dummy_token", nil
//}
