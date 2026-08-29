package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/noirbyss/worktrition-app/backend/user-service/internal/domain"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/password"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/token"
)

type AuthService struct {
	users         domain.UserRepository
	refreshTokens domain.RefreshTokenRepository
	tokens        *token.Service
}

func NewAuthService(
	users domain.UserRepository,
	refreshTokens domain.RefreshTokenRepository,
	tokens *token.Service,
) *AuthService {
	return &AuthService{
		users:         users,
		refreshTokens: refreshTokens,
		tokens:        tokens,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	name, email, plainPassword, birthDate string,
) (*domain.AuthSession, error) {
	if err := domain.ValidateCreateUser(name, email, plainPassword, birthDate); err != nil {
		return nil, err
	}

	parsedBirthDate, err := domain.ParseBirthDate(birthDate)
	if err != nil {
		return nil, err
	}

	passwordHash, err := password.Hash(plainPassword)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		Name:         strings.TrimSpace(name),
		Email:        strings.TrimSpace(email),
		PasswordHash: passwordHash,
		BirthDate:    parsedBirthDate,
	}

	if _, err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) Login(
	ctx context.Context,
	email, plainPassword string,
) (*domain.AuthSession, error) {
	email = strings.TrimSpace(email)

	if err := domain.ValidateCredentials(email, plainPassword); err != nil {
		return nil, err
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, err
	}

	if !password.Compare(user.PasswordHash, plainPassword) {
		return nil, domain.ErrInvalidCredentials
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, rawRefreshToken string) (*domain.AuthSession, error) {
	rawRefreshToken = strings.TrimSpace(rawRefreshToken)
	if rawRefreshToken == "" {
		return nil, domain.NewValidationError("refresh_token", "is required")
	}

	tokenHash := s.tokens.HashRefreshToken(rawRefreshToken)
	storedToken, err := s.refreshTokens.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if storedToken.RevokedAt != nil {
		return nil, domain.ErrTokenRevoked
	}
	if !storedToken.ExpiresAt.After(now) {
		return nil, domain.ErrTokenExpired
	}

	user, err := s.users.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokens.RevokeByHash(ctx, tokenHash); err != nil {
		return nil, err
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	rawRefreshToken = strings.TrimSpace(rawRefreshToken)
	if rawRefreshToken == "" {
		return domain.NewValidationError("refresh_token", "is required")
	}

	return s.refreshTokens.RevokeByHash(ctx, s.tokens.HashRefreshToken(rawRefreshToken))
}

func (s *AuthService) issueSession(ctx context.Context, user *domain.User) (*domain.AuthSession, error) {
	if user == nil {
		return nil, domain.NewValidationError("user", "is required")
	}

	now := time.Now()

	accessToken, err := s.tokens.NewAccessToken(user.ID, now)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := s.tokens.NewRefreshToken(now)
	if err != nil {
		return nil, err
	}

	storedToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshToken.Hash,
		ExpiresAt: refreshToken.ExpiresAt,
	}
	if err := s.refreshTokens.Create(ctx, storedToken); err != nil {
		return nil, err
	}

	return &domain.AuthSession{
		User:                  user,
		AccessToken:           accessToken.Value,
		RefreshToken:          refreshToken.Value,
		AccessTokenExpiresAt:  accessToken.ExpiresAt,
		RefreshTokenExpiresAt: refreshToken.ExpiresAt,
	}, nil
}
