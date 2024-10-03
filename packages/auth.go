package packages

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrValidation     = errors.New("validation error")
	ErrInternalServer = errors.New("internal server error")
)

type AuthResponse struct {
	AccessToken string
}

type AuthService interface {
	Register(ctx context.Context, input RegistrationInput) (*AuthResponse, error)
}

type RegistrationInput struct {
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type authService struct {
	// Add dependencies here, e.g., database, logger, etc.
}

func NewAuthService( /* Add necessary dependencies */ ) AuthService {
	return &authService{
		// Initialize dependencies
	}
}

func (s *authService) Register(ctx context.Context, input RegistrationInput) (*AuthResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	input.Sanitize()

	// Here you would typically:
	// 1. Check if the user already exists
	// 2. Hash the password
	// 3. Store the user in the database
	// 4. Generate and return an access token

	// For demonstration purposes:
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to hash password", ErrInternalServer)
	}

	// TODO: Implement user storage and token generation
	_ = hashedPassword // Placeholder to use hashedPassword

	return &AuthResponse{
		AccessToken: "dummy_token", // Replace with actual token generation
	}, nil
}

func (in *RegistrationInput) Sanitize() {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.Username = strings.TrimSpace(in.Username)
}

func (in *RegistrationInput) Validate() error {
	var errors []string

	if err := validateEmail(in.Email); err != nil {
		errors = append(errors, fmt.Sprintf("Email: %v", err))
	}

	if err := validateUsername(in.Username); err != nil {
		errors = append(errors, fmt.Sprintf("Username: %v", err))
	}

	if err := validatePassword(in.Password); err != nil {
		errors = append(errors, fmt.Sprintf("Password: %v", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", ErrValidation, strings.Join(errors, "; "))
	}

	return nil
}

func validateEmail(email string) error {
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
