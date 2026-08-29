package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) (string, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	SetProfileCompleted(ctx context.Context, id string, completed bool) error
}

type ProfileRepository interface {
	Save(ctx context.Context, profile *Profile) error
	GetByUserID(ctx context.Context, userID string) (*Profile, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeByHash(ctx context.Context, tokenHash string) error
	RevokeAllByUserID(ctx context.Context, userID string) error
}
