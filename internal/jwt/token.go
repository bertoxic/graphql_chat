package jwt

import (
	"context"
	"time"
)

var (
	AccessTokenLifeTime  = time.Minute * 15
	RefreshTokenLifeTime = time.Hour * 24 * 8
)

type RefreshToken struct {
	ID         string
	Sub        string
	Name       string
	LastUsedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreateRefreshTokenParams struct {
	UserID string
	Name   string
}
type RefreshTokenRepo interface {
	Create(ctx context.Context, createRefreshParams CreateRefreshTokenParams) (RefreshToken, error)
}
