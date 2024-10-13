package models

import (
	"errors"
	"fmt"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	Floatmap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

type InputDetails interface {
	Sanitize()
	Validate() error
}
type UserDetails struct {
	ID       string `bson:"id" json:"id"`
	UserName string `bson:"user_name" json:"user_name"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}
type RegistrationInput struct {
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type LoginInput struct {
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
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
		return fmt.Errorf("%w: %s", errorx.ErrValidation, strings.Join(errors, "; "))
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
