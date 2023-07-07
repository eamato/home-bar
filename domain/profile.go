package domain

import "context"

const (
	ProfilesTableName = "profiles"
)

type Profile struct {
	ID       int64  `json:"id" db:"id"`
	UserID   int64  `json:"user_id" db:"user_id"`
	Nickname string `json:"password" db:"nickname"`
}

type ProfileRepository interface {
	CreateProfile(context.Context, *Profile) error
	GetProfileByUserID(context.Context, int64) (*Profile, error)
}

type ProfileUsecase interface {
	GetProfile(context.Context, int64) (*Profile, *UsecaseError)
}
