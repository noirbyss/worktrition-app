package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/noirbyss/worktrition-app/backend/user-service/internal/domain"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/password"
)

type AuthService struct {
	users domain.UserRepository
}

func NewAuthService(users domain.UserRepository) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) CreateUser(
	ctx context.Context,
	name, email, plainPassword, birthDate string,
) (*domain.User, error) {
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

	return user, nil
}

func (s *AuthService) VerifyCredentials(
	ctx context.Context,
	email, plainPassword string,
) (*domain.User, error) {
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

	return user, nil
}
