package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/config"
	"github.com/bertoxic/graphqlChat/internal/utils"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var passwordCost = config.PasswordCost

type AuthService interface {
	Register(ctx context.Context, input RegistrationInput) (*AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (*AuthResponse, error)
}

type authService struct {
	// Add dependencies here, e.g., database, logger, etc.
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		// Initialize dependencies
	}
}

func (s *authService) Register(ctx context.Context, input RegistrationInput) (*AuthResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrValidation, err)
	}

	input.Sanitize()

	// Here you would typically:
	// 1. Check if the user already exists
	if _, err := s.userRepo.GetUserByEmail(ctx, input.Email); err != nil {
		return nil, fmt.Errorf("%w %v", err, utils.ErrUserExist)
	}
	// 2. Hash the password
	// 3. Store the user in the database
	// 4. Generate and return an access token

	// For demonstration purposes:
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), passwordCost)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to hash password", utils.ErrInternalServer)
	}

	// TODO: Implement user storage and token generation
	_ = hashedPassword // Placeholder to use hashedPassword
	user := &RegistrationInput{
		Email:    input.Email,
		Username: input.Username,
		Password: string(hashedPassword),
	}
	userDetails := UserDetails{}
	userDetails, err = s.userRepo.CreateUser(ctx, *user)
	if err != nil {
		return nil, err
	}
	// Generate access token (implementation details omitted)
	accessToken, err := generateAccessToken(user)
	return &AuthResponse{
		AccessToken: accessToken, // Replace with actual token generation
		User:        userDetails,
	}, nil
}
func (s *authService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrValidation, err)
	}

	input.Sanitize()

	// Check if the user exists
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return nil, fmt.Errorf("%w, %w: user not found", utils.ErrAuthentication, utils.ErrUserNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// TODO: Implement password verification here
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), passwordCost)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, fmt.Errorf("%w , %w, %w", err, utils.ErrAuthentication, utils.ErrInvalidCredentials)
	}
	// Generate access token
	accessToken, err := generateAccessToken(&input)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w, %w", err, utils.ErrInternalServer)
	}

	return &AuthResponse{
		AccessToken: accessToken,
		User:        UserDetails{Email: input.Email},
	}, nil
}
func verifyPassword(password string) (bool, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return false, utils.ErrValidation
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
	var errors []string

	if err := ValidateEmail(in.Email); err != nil {
		errors = append(errors, fmt.Sprintf("Email: %v", err))
	}

	if err := validateUsername(in.Username); err != nil {
		errors = append(errors, fmt.Sprintf("Username: %v", err))
	}

	if err := validatePassword(in.Password); err != nil {
		errors = append(errors, fmt.Sprintf("Password: %v", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(errors, "; "))
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
		return fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(errors, "; "))
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
