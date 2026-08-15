package auth

import (
	"ditdah/internal/features/user"
	jwtToken "ditdah/pkg/jwt"
	password "ditdah/pkg/jwt/crypto"
	"log/slog"

	// "crypto/rand"
	// "strings"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCase struct {
	users     user.UserRepository
	jwtSecret string
}

type LoginResult struct {
	Token        string
	RefreshToken string
}

func NewAuthUseCase(users user.UserRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		users:     users,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthUseCase) RegisterUseCase(ctx context.Context, in user.UserRegisterDTO) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	passwordHash, err := password.HashPassword(in.Password)
	slog.Info(passwordHash)
	if err != nil {
		slog.Error("Error hashing password", err)
		return err
	}

	return s.users.Register(ctx, in.Username, passwordHash)
}

func (s *AuthUseCase) Login(ctx context.Context, in LoginInput) (LoginResult, error) {
	in.Username = strings.TrimSpace(in.Username)
	if in.Username == "" || in.Password == "" {
		return LoginResult{}, errors.New("Invalid login or password")
	}

	jwtMng := jwtToken.NewTokenManager(s.jwtSecret)

	// if s.checkLoginRateLimit(ctx, in.Login) {
	// 	return LoginResult{}, ErrTooManyLoginAttempts
	// }

	user, err := s.users.GetUserByLogin(ctx, in.Username)
	if err != nil {
		return LoginResult{}, err
	}

	if user == nil {
		return LoginResult{}, errors.New("User is not exists")
	}

	if !password.CheckPasswordHash(in.Password, user.Password) {
		return LoginResult{}, errors.New("Invalid password")
	}

	accessToken, err := jwtMng.GenerateAccessToken(user.Id, in.Username)
	if err != nil {
		return LoginResult{}, err
	}

	refreshToken, err := jwtMng.GenerateRefreshToken()
	if err != nil {
		return LoginResult{}, err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	if err := s.users.SaveRefreshToken(ctx, user.Id, refreshToken, expiresAt); err != nil {
		return LoginResult{}, err
	}

	// s.resetLoginRateLimit(ctx, in.Username)

	return LoginResult{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthUseCase) Refresh(ctx context.Context, refreshToken string) (LoginResult, error) {
	session, err := s.users.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return LoginResult{}, err
	}

	if session == nil {
		return LoginResult{}, errors.New("invalid refresh token")
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.users.DeleteRefreshToken(ctx, refreshToken)
		return LoginResult{}, errors.New("refresh token expired")
	}

	user, err := s.users.GetFullUserById(ctx, session.UserID)
	if err != nil {
		return LoginResult{}, err
	}

	if user == nil {
		return LoginResult{}, errors.New("user not found")
	}

	jwtMng := jwtToken.NewTokenManager(s.jwtSecret)

	token, err := jwtMng.GenerateRefreshToken()
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: token}, nil
}

func (s *AuthUseCase) Logout(accessToken string) error {
	if accessToken == "" {
		return nil
	}

	claims := &jwtToken.Claims{}
	_, _, err := jwt.NewParser().ParseUnverified(accessToken, claims)
	if err != nil {
		return nil
	}

	var ttl time.Duration
	if claims.ExpiresAt != nil {
		ttl = time.Until(claims.ExpiresAt.Time)
	}
	if ttl <= 0 {
		return nil
	}

	return nil
}