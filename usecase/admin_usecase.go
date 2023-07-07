package usecase

import (
	"context"
	"home-bar/domain"
)

const (
	defaultTake = 1
	defaultSkip = 0
)

type adminUsecase struct {
	adminRepository domain.AdminRepository
	roleRepository  domain.RoleRepository
}

func NewAdminUsecase(adminRepository domain.AdminRepository, roleRepository domain.RoleRepository) domain.AdminUsecase {
	return &adminUsecase{
		adminRepository: adminRepository,
		roleRepository:  roleRepository,
	}
}

func (au *adminUsecase) GetUsers(
	ctx context.Context,
	usersListRequest domain.UsersListRequest,
) (*domain.UsersListResponse, *domain.UsecaseError) {
	take := defaultTake
	skip := defaultSkip

	if usersListRequest.PaginationRequest.Take > 0 {
		take = usersListRequest.PaginationRequest.Take
	}

	if usersListRequest.PaginationRequest.Skip > 0 {
		skip = usersListRequest.PaginationRequest.Skip
	}

	res, err := au.adminRepository.GetUsersFullInfo(ctx, take, skip)
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	return &domain.UsersListResponse{
		Users: res,
	}, nil
}

func (au *adminUsecase) DeleteUser(
	ctx context.Context,
	request domain.UserDeleteRequest,
) (*domain.UserDeleteResponse, *domain.UsecaseError) {
	isSuccessful, err := au.adminRepository.DeleteUser(ctx, request.UserID)
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	return &domain.UserDeleteResponse{
		UserID:  request.UserID,
		Deleted: isSuccessful,
	}, nil
}

func (au *adminUsecase) AssignRole(
	ctx context.Context,
	request domain.RoleAssignRequest,
) (*domain.RoleAssignResponse, *domain.UsecaseError) {
	id, err := au.adminRepository.UpdateUserRole(ctx, request.UserID, request.RoleID)
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	role, err := au.roleRepository.GetRole(ctx, request.UserID)
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	return &domain.RoleAssignResponse{
		ID:     id,
		UserID: request.UserID,
		Role:   *role,
	}, nil
}
