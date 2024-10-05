package auth

import (
	"context"
	"errors"
	"github.com/bertoxic/graphqlChat/internal/utils"
	"github.com/bertoxic/graphqlChat/test"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	// mocks "github.com/bertoxic/graphqlChat/mocks/packages"
)

//nolint:funlen

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)
	tests := []struct {
		name            string
		input           RegistrationInput
		setupMock       func(*MockUserRepository)
		wantErr         bool
		errType         error
		expectedCalled  bool
		expectedErr     error
		expectedCalls   map[string][]interface{}
		unexpectedCalls []string
	}{
		{
			name: "Valid registration",
			input: RegistrationInput{
				Email:    "test@example.com",
				Username: "validuser",
				Password: "V@lidP@ssw0rd",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("GetUserByEmail",
					mock.Anything,
					"test@example.com",
				).Return(UserDetails{ID: "123"}, nil)
				m.On("CreateUser", mock.Anything, mock.AnythingOfType("RegistrationInput")).Return(UserDetails{ID: "1234"}, nil)
			},
			wantErr:        false,
			expectedCalled: true,
			expectedCalls: map[string][]interface{}{
				"GetUserByEmail": {mock.Anything, "test@example.com"},
				"CreateUser":     {mock.Anything, mock.AnythingOfType("RegistrationInput")}},
		},
		{
			name: "Invalid email",
			input: RegistrationInput{
				Email:    "invalid-email",
				Username: "validuser",
				Password: "V@lidP@ssw0rd",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "test@example.com", "validuser").Return(false, nil)
				m.On("CreateUser", mock.Anything, mock.AnythingOfType("User")).Return(nil)

			},
			wantErr: true,
			errType: utils.ErrValidation,
		},
		{
			name: "user already exist",
			input: RegistrationInput{
				Email:    "example@gmail.com",
				Username: "validateUsername",
				Password: "validP@ssw0rd",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything,
					"example@gmail.com").Return(
					UserDetails{}, errors.New("user already exist"))
			},
			wantErr:        true,
			errType:        utils.ErrUserExist,
			expectedCalled: true,
			expectedCalls: map[string][]interface{}{
				"GetUserByEmail": {mock.Anything, "example@gmail.com"},
			},
		},
		{
			name: "Invalid username - too short",
			input: RegistrationInput{
				Email:    "test@example.com",
				Username: "ab",
				Password: "V@lidP@ssw0rd",
			},
			wantErr: true,
			errType: utils.ErrValidation,
		},
		{
			name: "Invalid username - invalid characters",
			input: RegistrationInput{
				Email:    "test@example.com",
				Username: "invalid user!",
				Password: "V@lidP@ssw0rd",
			},
			wantErr: true,
			errType: utils.ErrValidation,
		},
		{
			name: "Invalid password - too short",
			input: RegistrationInput{
				Email:    "test@example.com",
				Username: "validuser",
				Password: "Short1!",
			},
			wantErr: true,
			errType: utils.ErrValidation,
		},
		{
			name: "Invalid password - missing requirements",
			input: RegistrationInput{
				Email:    "test@example.com",
				Username: "validuser",
				Password: "onlylowercase",
			},
			wantErr: true,
			errType: utils.ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil // Clear previous mock expectations
			mockRepo.Calls = nil         // Clear previous mock calls

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			resp, err := service.Register(context.Background(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errType.Error())
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp.AccessToken)
			}

			if tt.expectedCalled {
				for methodName, args := range tt.expectedCalls {
					mockRepo.AssertCalled(t, methodName, args...)
				}
				mockRepo.AssertExpectations(t)
			} else {
				for _, call := range tt.unexpectedCalls {
					mockRepo.AssertNotCalled(t, call)
				}
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
			assert.NoError(t, ValidateEmail(email))
		})
	}

	for _, email := range invalidEmails {
		t.Run(email, func(t *testing.T) {
			assert.Error(t, ValidateEmail(email))
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

func TestAuthService_Register2(t *testing.T) {
	ctx := context.Background()
	validInput := RegistrationInput{
		Email:    "someemail@email.com",
		Username: "AValidUser",
		Password: "passW@and45",
	}

	t.Run("can register", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		service := NewAuthService(userRepo)
		userRepo.On("GetUserByEmail", mock.Anything, mock.Anything, mock.Anything).Return(UserDetails{}, nil)
		userRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("RegistrationInput")).Return(UserDetails{
			ID:       "1234",
			UserName: "name",
			Email:    "some@fm.com",
		}, nil)

		resp, err := service.Register(ctx, validInput)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "1234", resp.User.ID)
		require.Equal(t, "name", resp.User.UserName)
		require.Equal(t, "some@fm.com", resp.User.Email)
		require.NotEmpty(t, resp.AccessToken)
	})
	t.Run("can register", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		service := NewAuthService(userRepo)
		userRepo.On("GetUserByEmail", mock.Anything, mock.Anything, mock.Anything).Return(UserDetails{}, utils.ErrUserExist)
		_, err := service.Register(ctx, validInput)
		require.Error(t, err)
		require.ErrorIs(t, err, utils.ErrUserExist)
		userRepo.AssertCalled(t, "GetUserByEmail", mock.Anything, mock.Anything, mock.Anything)
		userRepo.AssertNotCalled(t, "CreateUser")
		require.NotEmpty(t, err)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)
	ctx := context.Background()
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password3@"), config.PasswordCost)
	tests := []struct {
		name            string
		input           LoginInput
		expectedError   error
		setupMock       func(*MockUserRepository)
		wantErr         bool
		expectedCalls   map[string][]interface{}
		unexpectedCalls []string
	}{
		{
			name: "valid inputs",
			input: LoginInput{
				Email:    "validmail@gmail.com",
				Password: "Password3@",
			},
			wantErr: false,
			setupMock: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "validmail@gmail.com").Return(UserDetails{Password: test.FakePassWord}, nil)
			},
			expectedCalls: map[string][]interface{}{
				"GetUserByEmail": {mock.Anything, "validmail@gmail.com"},
			},
		},
		{
			name: "User does not exist",
			input: LoginInput{
				Email:    "validmail@gmail.com",
				Password: "AvEry1@va#lid",
			},
			wantErr: true,
			setupMock: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "validmail@gmail.com").Return(UserDetails{}, utils.ErrUserNotFound)
			},
			expectedCalls: map[string][]interface{}{
				"GetUserByEmail": {mock.Anything, "validmail@gmail.com"},
			},
			expectedError: utils.ErrUserNotFound,
		},
		{
			name: "wrong password",
			input: LoginInput{
				Email:    "validmail@gmail.com",
				Password: "wrongPassWord3@",
			},
			wantErr: true,
			setupMock: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "validmail@gmail.com").Return(UserDetails{Password: test.FakePassWord}, nil)
			},
			expectedCalls: map[string][]interface{}{
				"GetUserByEmail": {mock.Anything, "validmail@gmail.com"},
			},
			expectedError: utils.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil // Clear previous mock expectations
			mockRepo.Calls = nil
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}
			_, err := service.Login(ctx, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}

			for methodName, args := range tt.expectedCalls {
				mockRepo.AssertCalled(t, methodName, args...)
			}

			for _, methodName := range tt.unexpectedCalls {
				mockRepo.AssertNotCalled(t, methodName)
			}
		})
	}
}
