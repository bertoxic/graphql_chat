package middlewares

import (
	"context"
	"github.com/bertoxic/graphqlChat/internal/auth"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"net/http"
)

type contextKey string

var (
	contextAuthIDKey contextKey = "currentUserID"
)

func AuthMiddleWare(authTokenService auth.TokenService) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := authTokenService.ParseTokenFromRequest(ctx, r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx = PutUserIDINContext(ctx, token.Sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	if ctx.Value(contextAuthIDKey) == nil {
		return "", errorx.New(errorx.ErrCodeNoUserIdInContext, "no user id in context", nil)
	}
	userID, ok := ctx.Value(contextAuthIDKey).(string)
	if !ok {
		return "", errorx.New(errorx.ErrCodeNoUserIdInContext, "no user id in context", nil)
	}
	return userID, nil
}

func PutUserIDINContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextAuthIDKey, userID)
}
