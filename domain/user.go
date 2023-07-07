package domain

import "context"

const (
	UsersTableName = "users"
)

type User struct {
	ID           int64  `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email" db:"email"`
	Password     string `json:"password" db:"password"`
	RefreshToken string `db:"refresh_token"`
}

type UserRepository interface {
	CreateUser(context.Context, *User) (int64, error)
	GetUserByName(context.Context, string) (*User, error)
	UpdateUser(context.Context, *User) (int64, error)
	GetUserById(context.Context, int64) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
	GetUserByUsernameOrEmail(context.Context, string, string) (*User, error)
}
