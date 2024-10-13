package jwt

import (
	"context"
	"github.com/bertoxic/graphqlChat/internal/auth"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"github.com/lestrrat-go/jwx/jwa"
	jwtGo "github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"time"
)

var SignatureType = jwa.HS256

type TokenService struct {
	config *config.AppConfig
}

func NewTokenService(config *config.AppConfig) *TokenService {
	return &TokenService{config: config}
}

func (ts *TokenService) ParseTokenFromRequest(ctx context.Context, r *http.Request) (auth.AuthToken, error) {
	token, err := jwtGo.ParseRequest(r,
		jwtGo.WithValidate(true),
		jwtGo.WithIssuer(ts.config.JWT.Issuer),
		jwtGo.WithVerify(SignatureType, ts.config.JWT.Secret),
	)
	if err != nil {
		err = errorx.New(errorx.ErrCodeInvalidToken, "invalid accessToken", err)
		return auth.AuthToken{}, err
	}
	return buildToken(token), nil
}

func buildToken(token jwtGo.Token) auth.AuthToken {
	return auth.AuthToken{
		TokenID: token.JwtID(),
		Sub:     token.Subject(),
	}
}

func (ts *TokenService) ParseToken(ctx context.Context, payload string) (auth.AuthToken, error) {
	token, err := jwtGo.Parse(
		[]byte(payload),
		jwtGo.WithValidate(true),
		jwtGo.WithIssuer(ts.config.JWT.Issuer),
		jwtGo.WithVerify(SignatureType, ts.config.JWT.Secret),
	)
	if err != nil {
		err = errorx.New(errorx.ErrCodeInvalidToken, "invalid accessToken", err)
		return auth.AuthToken{}, err
	}
	return buildToken(token), nil
}

func (ts *TokenService) CreateAccessToken(ctx context.Context, user models.UserDetails) (string, error) {
	t := jwtGo.New()
	if err := setDefaultToken(t, user, AccessTokenLifeTime, ts.config); err != nil {
		return "", errorx.New(errorx.ErrCodeValidation, "cannot set default token", err)
	}
	signedToken, err := jwtGo.Sign(t, SignatureType, ts.config.JWT.Secret)
	if err != nil {
		return "", errorx.New(errorx.ErrCodeValidation, "cannot sign default token", err)
	}
	return string(signedToken), nil
}
func (ts *TokenService) CreateRefreshToken(ctx context.Context, user models.UserDetails, tokenID string) (string, error) {
	t := jwtGo.New()
	if err := setDefaultToken(t, user, RefreshTokenLifeTime, ts.config); err != nil {
		return "", errorx.New(errorx.ErrCodeValidation, "cannot set default token", err)
	}
	err := t.Set(jwtGo.JwtIDKey, tokenID)
	if err != nil {
		return "", errorx.New(errorx.ErrCodeValidation, "cannot sign default token", err)
	}
	signedToken, err := jwtGo.Sign(t, SignatureType, ts.config.JWT.Secret)
	if err != nil {
		return "", errorx.New(errorx.ErrCodeValidation, "cannot sign default token", err)
	}
	return string(signedToken), nil
}

func setDefaultToken(token jwtGo.Token, user models.UserDetails, lifetime time.Duration, config *config.AppConfig) error {
	if err := token.Set(jwtGo.SubjectKey, user.ID); err != nil {
		return errorx.New(errorx.ErrCodeValidation, "error setting jwt sub", err)
	}
	if err := token.Set(jwtGo.IssuerKey, config.JWT.Issuer); err != nil {
		return errorx.New(errorx.ErrCodeValidation, "error setting jwt sub", err)
	}
	if err := token.Set(jwtGo.IssuedAtKey, time.Now().Unix()); err != nil {
		return errorx.New(errorx.ErrCodeValidation, "error setting jwt sub", err)
	}
	if err := token.Set(jwtGo.ExpirationKey, time.Now().Add(lifetime).Unix()); err != nil {
		return errorx.New(errorx.ErrCodeValidation, "error setting jwt sub", err)
	}
	return nil
}
