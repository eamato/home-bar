package domain

import "context"

type UsersListRequest struct {
	PaginationRequest PaginationRequest `form:"pagination"`
}

type UsersListResponseProfile struct {
	ID       int64  `json:"id" db:"profile_id"`
	UserID   int64  `json:"user_id" db:"user_id"`
	Nickname string `json:"password" db:"nickname"`
}

type UsersListResponseUser struct {
	ID                       int64  `json:"id" db:"id"`
	Username                 string `json:"username" db:"username"`
	Email                    string `json:"email" db:"email"`
	Password                 string `json:"password" db:"password"`
	RefreshToken             string `json:"refresh_token" db:"refresh_token"`
	UsersListResponseProfile `json:"profile"`
	Role                     string `json:"role" db:"role"`
}

type UsersListResponse struct {
	Users []*UsersListResponseUser `json:"users"`
}

type UserDeleteRequest struct {
	UserID int64 `form:"user_id" binding:"required"`
}

type UserDeleteResponse struct {
	UserID  int64 `json:"user_id"`
	Deleted bool  `json:"deleted"`
}

type RoleAssignRequest struct {
	UserID int64 `form:"user_id" binding:"required"`
	RoleID int64 `form:"role_id" binding:"required"`
}

type RoleAssignResponse struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Role   Role  `json:"role"`
}

type AdminUsecase interface {
	GetUsers(context.Context, UsersListRequest) (*UsersListResponse, *UsecaseError)
	DeleteUser(context.Context, UserDeleteRequest) (*UserDeleteResponse, *UsecaseError)
	AssignRole(context.Context, RoleAssignRequest) (*RoleAssignResponse, *UsecaseError)
}

type AdminRepository interface {
	GetUsersFullInfo(context.Context, int, int) ([]*UsersListResponseUser, error)
	DeleteUser(context.Context, int64) (bool, error)
	UpdateUserRole(context.Context, int64, int64) (int64, error)
}
