package packages

import (
	"context"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {
	service := NewAuthService()

	tests := []struct {
		name    string
		input   RegistrationInput
		wantErr bool
		errType error
	}{
		{
			name: "Valid registration",
			input: RegistrationInput{
				Email:    "test@example.com",
				Username: "validuser",
				Password: "V@lidP@ssw0rd",
			},
			wantErr: false,
		},
		{
			name: "Invalid email",
			input: RegistrationInput{
				Email:    "invalid-email",
				Username: "validuser",
				Password: "V@lidP@ssw0rd",
			},
			wantErr: true,
			errType: ErrValidation,
		},
		{
			name: "Invalid username - too short",
			input: RegistrationInput{
				Email:    "test@example.com",
				Username: "ab",
				Password: "V@lidP@ssw0rd",
			},
			wantErr: true,
			errType: ErrValidation,
		},
		{
			name: "Invalid username - invalid characters",
			input: RegistrationInput{
				Email:    "test@example.com",
				Username: "invalid user!",
				Password: "V@lidP@ssw0rd",
			},
			wantErr: true,
			errType: ErrValidation,
		},
		{
			name: "Invalid password - too short",
			input: RegistrationInput{
				Email:    "test@example.com",
				Username: "validuser",
				Password: "Short1!",
			},
			wantErr: true,
			errType: ErrValidation,
		},
		{
			name: "Invalid password - missing requirements",
			input: RegistrationInput{
				Email:    "test@example.com",
				Username: "validuser",
				Password: "onlylowercase",
			},
			wantErr: true,
			errType: ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Register(context.Background(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.errType)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp.AccessToken)
			}
		})
	}
}

func TestRegistrationInput_Sanitize(t *testing.T) {
	input := RegistrationInput{
		Email:    " TEST@EXAMPLE.COM ",
		Username: " testUser ",
		Password: "password123",
	}

	input.Sanitize()

	assert.Equal(t, "test@example.com", input.Email)
	assert.Equal(t, "testUser", input.Username)
	assert.Equal(t, "password123", input.Password, "Password should not be modified")
}

func TestValidateEmail(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name+tag+sorting@example.com",
		"x@example.com",
	}

	invalidEmails := []string{
		"invalid-email",
		"@missing-local.org",
		"missing-at-sign.net",
		"missing-domain@.com",
	}

	for _, email := range validEmails {
		t.Run(email, func(t *testing.T) {
			assert.NoError(t, validateEmail(email))
		})
	}

	for _, email := range invalidEmails {
		t.Run(email, func(t *testing.T) {
			assert.Error(t, validateEmail(email))
		})
	}
}

func TestValidateUsername(t *testing.T) {
	validUsernames := []string{
		"validuser",
		"user123",
		"user_name",
		"user-name",
	}

	invalidUsernames := []string{
		"ab",                          // Too short
		"thisusernameiswaytooooolong", // Too long
		"invalid user",                // Contains space
		"invalid@user",                // Contains @
	}

	for _, username := range validUsernames {
		t.Run(username, func(t *testing.T) {
			assert.NoError(t, validateUsername(username))
		})
	}

	for _, username := range invalidUsernames {
		t.Run(username, func(t *testing.T) {
			assert.Error(t, validateUsername(username))
		})
	}
}

func TestValidatePassword(t *testing.T) {
	validPasswords := []string{
		"V@lidP@ssw0rd",
		"Str0ng!Pass",
		"C0mpl3x@Pass",
	}

	invalidPasswords := []string{
		"short",             // Too short
		"onlylowercase",     // Missing uppercase, number, and special char
		"ONLYUPPERCASE",     // Missing lowercase, number, and special char
		"onlyalphanumeric1", // Missing uppercase and special char
	}

	for _, password := range validPasswords {
		t.Run(password, func(t *testing.T) {
			assert.NoError(t, validatePassword(password))
		})
	}

	for _, password := range invalidPasswords {
		t.Run(password, func(t *testing.T) {
			assert.Error(t, validatePassword(password))
		})
	}
}
