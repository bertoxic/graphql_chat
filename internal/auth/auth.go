package auth

import (
	"context"
	"errors"
	"fmt"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var passwordCost = config.PasswordCost

type AuthService interface {
	Register(ctx context.Context, input models.RegistrationInput) (*AuthResponse, error)
	Login(ctx context.Context, input models.LoginInput) (*AuthResponse, error)
}

type authService struct {
	userRepo UserRepository
	Service  TokenService
}
type AuthToken struct {
	TokenID string
	Sub     string
}

type TokenService interface {
	ParseTokenFromRequest(ctx context.Context, r *http.Request) (AuthToken, error)
	ParseToken(ctx context.Context, payload string) (AuthToken, error)
	CreateAccessToken(ctx context.Context, user models.UserDetails) (string, error)
	CreateRefreshToken(ctx context.Context, user models.UserDetails, tokenID string) (string, error)
}

func NewAuthService(userRepo UserRepository, service TokenService) AuthService {
	return &authService{
		userRepo: userRepo,
		Service:  service,
	}
}

func (s *authService) Register(ctx context.Context, input models.RegistrationInput) (*AuthResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errorx.New(errorx.ErrCodeValidation, "registration failed", err), err)
	}

	input.Sanitize()

	if _, err := s.userRepo.GetUserByEmail(ctx, input.Email); err != nil {
		if errors.As(err, &errorx.ErrNotFound) {

		} else {
			return nil, fmt.Errorf("%w %v", err, errorx.ErrNotFound)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), passwordCost)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to hash password", errorx.ErrInternal)
	}

	// TODO: Implement user storage and token generation
	_ = hashedPassword
	user := &models.RegistrationInput{
		Email:    input.Email,
		Username: input.Username,
		Password: string(hashedPassword),
	}

	userDetails := &models.UserDetails{}
	userDetails, err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.Service.CreateAccessToken(ctx, *userDetails)
	return &AuthResponse{
		AccessToken: accessToken,
		User:        *userDetails,
	}, nil
}
func (s *authService) Login(ctx context.Context, input models.LoginInput) (*AuthResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errorx.ErrValidation, err)
	}
	input.Sanitize()

	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, errorx.ErrNotFound) {
			return nil, fmt.Errorf("%w, %w: user not found", errorx.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// TODO: Implement password verification here

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, fmt.Errorf("%w , %w, %w", err, errorx.ErrInvalidCredentials)
	}
	accessToken, err := s.Service.CreateAccessToken(ctx, *user)
	if err != nil {
		return nil, err
	}
	//accessToken, err := generateAccessToken(&input)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w, %w", err, errorx.ErrInternal)
	}

	return &AuthResponse{
		AccessToken: accessToken,
		User: models.UserDetails{
			ID:       user.ID,
			UserName: user.UserName,
			Email:    user.Email,
			Password: user.Password,
		},
	}, nil
}
func verifyPassword(password string) (bool, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return false, errorx.ErrValidation
	}
	return true, nil
}
func generateAccessToken(user InputDetails) (string, error) {
	user.Sanitize()
	if err := user.Validate(); err != nil {
		return "", err
	}
	return "dummy_token", nil
}

func (in *RegistrationInput) Sanitize() {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.Username = strings.TrimSpace(in.Username)
}
func (in *RegistrationInput) Validate() error {
	var errorList []string

	if err := ValidateEmail(in.Email); err != nil {
		errorList = append(errorList, fmt.Sprintf("Email: %v", err))
	}

	if err := validateUsername(in.Username); err != nil {
		errorList = append(errorList, fmt.Sprintf("Username: %v", err))
	}

	if err := validatePassword(in.Password); err != nil {
		errorList = append(errorList, fmt.Sprintf("Password: %v", err))
	}

	if len(errorList) > 0 {
		return fmt.Errorf("%w: %s", errorx.ErrValidation, strings.Join(errorList, "; "))
	}

	return nil
}
func (in *LoginInput) Sanitize() {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
}

func (in *LoginInput) Validate() error {
	var errors []string

	if err := ValidateEmail(in.Email); err != nil {
		errors = append(errors, fmt.Sprintf("Email: %v", err))
	}

	if err := validatePassword(in.Password); err != nil {
		errors = append(errors, fmt.Sprintf("Password: %v", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", errorx.ErrValidation, strings.Join(errors, "; "))
	}

	return nil
}

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email format")
	}
	return nil
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("must be between 3 and 20 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username) {
		return errors.New("can only contain alphanumeric characters, underscores, and hyphens")
	}
	return nil
}

func validatePassword(password string) error {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if !hasMinLen {
		return errors.New("must be at least 8 characters long")
	}
	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("must include at least one uppercase letter, one lowercase letter, one number, and one special character")
	}
	return nil
}
