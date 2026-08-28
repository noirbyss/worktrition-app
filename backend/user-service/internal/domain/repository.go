package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) (string, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

type ProfileRepository interface {
	Save(ctx context.Context, profile *Profile) error
	GetByUserID(ctx context.Context, userID string) (*Profile, error)
}
