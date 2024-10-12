package resolvers

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/bertoxic/graphqlChat/internal/auth"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthService auth.AuthService
	UserService auth.UserRepository
}

func NewResolver(authService auth.AuthService, userService auth.UserRepository) *Resolver {
	return &Resolver{
		AuthService: authService,
		UserService: userService,
	}
}

func buildBadRequestError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Err:       err,
		Message:   err.Error(),
		Path:      graphql.GetPath(ctx),
		Locations: nil,
		Extensions: map[string]interface{}{
			"code": http.StatusBadRequest,
		},
	}
}
