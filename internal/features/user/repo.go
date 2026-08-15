package user

import (
	"context"
	"time"
)

type UserRepository interface {
	// work with users
	GetAllUsers(ctx context.Context) ([]*UserEntity, error)
	GetFullUserById(ctx context.Context, userID int) (*UserEntity, error)
	GetUserByLogin(ctx context.Context, login string) (*UserEntity, error)

	// refresh system
	SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error

	// auth
	Register(ctx context.Context, username, passwordHash string) error
}