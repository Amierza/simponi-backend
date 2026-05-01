package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IUserService interface {
		CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
		GetUsers(ctx context.Context, req *response.PaginationRequest) (dto.UserPaginationResponse, error)
		GetUserByUserID(ctx context.Context, userID *uuid.UUID) (*dto.UserResponse, error)
		GetUserProfile(ctx context.Context) (*dto.UserResponse, error)
		UpdateUserByUserID(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
		UpdateUserStatusByUserID(ctx context.Context, req *dto.UpdateUserStatus) (*dto.UserResponse, error)
		UpdateUserProfile(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
		DeleteUserByUserID(ctx context.Context, userID *uuid.UUID) error
	}

	userService struct {
		userRepo   repository.IUserRepository
		roleRepo   repository.IRoleRepository
		logger     *zap.Logger
		jwtService jwt.IJWT
	}
)

func NewUserService(userRepo repository.IUserRepository, roleRepo repository.IRoleRepository, logger *zap.Logger, jwtService jwt.IJWT) *userService {
	return &userService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		logger:     logger,
		jwtService: jwtService,
	}
}

func mapToUserResponse(u *entity.User, r *entity.Role) *dto.UserResponse {

	return &dto.UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		ImageURL: u.ImageURL,
		Status:   u.Status,
		Role: dto.RoleResponse{
			ID:   r.ID,
			Name: r.Name,
		},
	}
}

func mapToProfileResponse(u *entity.User, r *entity.Role) *dto.UserResponse {
	permissions := []dto.PermissionResponse{}

	if r.RolePermissions != nil {
		for _, rp := range r.RolePermissions {
			permissions = append(permissions, dto.PermissionResponse{
				ID:       rp.Permission.ID,
				Name:     rp.Permission.Name,
				Endpoint: rp.Permission.Endpoint,
				Method:   rp.Permission.Method,
			})
		}
	}

	return &dto.UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		ImageURL: u.ImageURL,
		Role: dto.RoleResponse{
			ID:          r.ID,
			Name:        r.Name,
			Permissions: permissions,
		},
	}
}

func (us *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	_, found, err := us.userRepo.GetUserByEmail(ctx, nil, &req.Email)
	if err != nil {
		us.logger.Error("failed to get user by email", zap.String("email", req.Email), zap.Error(err))
		return nil, fmt.Errorf("failed to get user by email: %w", dto.ErrInternal)
	}
	if found {
		us.logger.Warn("user already exists", zap.String("email", req.Email))
		return nil, fmt.Errorf("user already exists: %w", dto.ErrAlreadyExists)
	}

	role, found, err := us.roleRepo.GetRoleByRoleID(ctx, nil, req.RoleID)
	if err != nil {
		us.logger.Error("failed to get role by id", zap.String("role_id", req.RoleID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get role by id: %w", dto.ErrInternal)
	}
	if !found {
		us.logger.Warn("role not found", zap.String("role_id", req.RoleID.String()))
		return nil, fmt.Errorf("role not found: %w", dto.ErrNotFound)
	}

	newID := uuid.New()
	newUser := &entity.User{
		ID:       newID,
		Name:     req.Name,
		Email:    req.Email,
		ImageURL: req.ImageURL,
		Password: req.Password,
		RoleID:   req.RoleID,
	}
	err = us.userRepo.CreateUser(ctx, nil, newUser)
	if err != nil {
		us.logger.Error("failed to create user", zap.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", dto.ErrInternal)
	}

	us.logger.Info("success to create user", zap.String("id", newUser.ID.String()))

	return mapToUserResponse(newUser, role), nil
}

func (us *userService) GetUsers(ctx context.Context, req *response.PaginationRequest) (dto.UserPaginationResponse, error) {
	datas, err := us.userRepo.GetUsers(ctx, nil, req)
	if err != nil {
		us.logger.Error("failed to get users", zap.Error(err))
		return dto.UserPaginationResponse{}, fmt.Errorf("failed to get users: %w", dto.ErrInternal)
	}

	us.logger.Info("success to get users", zap.Int64("count", datas.Count))

	var users []*dto.UserResponse
	for _, user := range datas.Users {
		users = append(users, mapToUserResponse(user, &user.Role))
	}

	return dto.UserPaginationResponse{
		Data:               users,
		PaginationResponse: datas.PaginationResponse,
	}, nil
}

func (us *userService) GetUserByUserID(ctx context.Context, userID *uuid.UUID) (*dto.UserResponse, error) {
	user, found, err := us.userRepo.GetUserByUserID(ctx, nil, userID)
	if err != nil {
		us.logger.Error("failed to get user by ID", zap.String("userID", userID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get user ID: %w", dto.ErrInternal)
	}
	if !found {
		us.logger.Warn("user not found", zap.String("userID", userID.String()))
		return nil, fmt.Errorf("user not found: %w", dto.ErrNotFound)
	}

	us.logger.Info("success to get user by id", zap.String("id", userID.String()))

	return mapToUserResponse(user, &user.Role), nil
}

func (us *userService) GetUserProfile(ctx context.Context) (*dto.UserResponse, error) {
	userIDString := ctx.Value("user_id").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		us.logger.Error("failed to parse user_id", zap.String("user_id", userIDString), zap.Error(err))
		return nil, fmt.Errorf("failed to create store: %w", dto.ErrInternal)
	}

	data, found, err := us.userRepo.GetUserByUserID(ctx, nil, &userID)
	if err != nil {
		us.logger.Error("failed to get user by id", zap.String("id", userID.String()), zap.Error((err)))
		return nil, fmt.Errorf("failed to get user by id: %w", dto.ErrInternal)
	}
	if !found {
		us.logger.Warn("user not found", zap.String("id", userID.String()))
		return nil, fmt.Errorf("user not found: %w", dto.ErrNotFound)
	}

	return mapToProfileResponse(data, &data.Role), nil
}

func (us *userService) UpdateUserByUserID(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, found, err := us.userRepo.GetUserByUserID(ctx, nil, &req.ID)
	if err != nil {
		us.logger.Error("failed to get user by ID", zap.String("userID", req.ID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get user ID: %w", dto.ErrInternal)
	}
	if !found {
		us.logger.Warn("user not found", zap.String("userID", req.ID.String()))
		return nil, fmt.Errorf("user not found: %w", dto.ErrNotFound)
	}

	// validate email
	if user.Email != req.Email {
		_, found, err = us.userRepo.GetUserByEmail(ctx, nil, &req.Email)
		if err != nil {
			us.logger.Error("failed to get user by email", zap.String("email", req.Email), zap.Error(err))
			return nil, fmt.Errorf("failed to get user by email: %w", dto.ErrInternal)
		}

		if found {
			us.logger.Warn("user email already exists", zap.String("email", req.Email))
			return nil, fmt.Errorf("user email already exists: %w", dto.ErrAlreadyExists)
		}
	}
	user.Email = req.Email

	// validate role
	if req.RoleID != nil {
		if user.RoleID != req.RoleID {
			role, found, err := us.roleRepo.GetRoleByRoleID(ctx, nil, req.RoleID)
			if err != nil {
				us.logger.Error("failed to get role by id", zap.String("role_id", req.RoleID.String()), zap.Error(err))
				return nil, fmt.Errorf("failed to get role by id: %w", dto.ErrInternal)
			}
			if !found {
				us.logger.Warn("role not found", zap.String("role_id", req.RoleID.String()))
				return nil, fmt.Errorf("role not found: %w", dto.ErrNotFound)
			}
			user.Role = *role
		}
		user.RoleID = req.RoleID
	}

	if req.ImageURL != nil {
		user.ImageURL = *req.ImageURL
	}
	user.Name = req.Name

	err = us.userRepo.UpdateUserByUserID(ctx, nil, user)
	if err != nil {
		us.logger.Error("failed to update user", zap.String("id", req.ID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to update user: %w", dto.ErrInternal)
	}

	return mapToUserResponse(user, &user.Role), nil
}

func (us *userService) UpdateUserStatusByUserID(ctx context.Context, req *dto.UpdateUserStatus) (*dto.UserResponse, error) {
	user, found, err := us.userRepo.GetUserByUserID(ctx, nil, &req.ID)
	if err != nil {
		us.logger.Error("failed to get user by ID", zap.String("userID", req.ID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get user ID: %w", dto.ErrInternal)
	}
	if !found {
		us.logger.Warn("user not found", zap.String("userID", req.ID.String()))
		return nil, fmt.Errorf("user not found: %w", dto.ErrNotFound)
	}

	// skip if same status
	if user.Status == req.Status {
		return mapToUserResponse(user, &user.Role), nil
	}

	// validate status
	if req.Status != "active" && req.Status != "inactive" {
		us.logger.Warn("invalid status", zap.String("status", req.Status))
		return nil, fmt.Errorf("invalid status: %w", dto.ErrBadRequest)
	}
	user.Status = req.Status
	log.Println(user.Status)

	err = us.userRepo.UpdateUserStatusByUserID(ctx, nil, user)
	if err != nil {
		us.logger.Error("failed to update user status", zap.String("status", req.Status), zap.Error(err))
		return nil, fmt.Errorf("failed to update user status: %w", dto.ErrInternal)
	}

	return mapToUserResponse(user, &user.Role), nil
}

func (us *userService) UpdateUserProfile(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	userIDString := ctx.Value("user_id").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		us.logger.Error("failed to parse user_id", zap.String("user_id", userIDString), zap.Error(err))
		return nil, fmt.Errorf("failed to create store: %w", dto.ErrInternal)
	}

	user, found, err := us.userRepo.GetUserByUserID(ctx, nil, &userID)
	if err != nil {
		us.logger.Error("failed to get user by ID", zap.String("userID", userID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get user ID: %w", dto.ErrInternal)
	}
	if !found {
		us.logger.Warn("user not found", zap.String("userID", userID.String()))
		return nil, fmt.Errorf("user not found: %w", dto.ErrNotFound)
	}

	// validate email
	if user.Email != req.Email {
		_, found, err = us.userRepo.GetUserByEmail(ctx, nil, &req.Email)
		if err != nil {
			us.logger.Error("failed to get user by email", zap.String("email", req.Email), zap.Error(err))
			return nil, fmt.Errorf("failed to get user by email: %w", dto.ErrInternal)
		}

		if found {
			us.logger.Warn("user email already exists", zap.String("email", req.Email))
			return nil, fmt.Errorf("user email already exists: %w", dto.ErrAlreadyExists)
		}
	}
	user.Email = req.Email

	// validate role
	if req.RoleID != nil {
		if user.RoleID != req.RoleID {
			role, found, err := us.roleRepo.GetRoleByRoleID(ctx, nil, req.RoleID)
			if err != nil {
				us.logger.Error("failed to get role by id", zap.String("role_id", req.RoleID.String()), zap.Error(err))
				return nil, fmt.Errorf("failed to get role by id: %w", dto.ErrInternal)
			}
			if !found {
				us.logger.Warn("role not found", zap.String("role_id", req.RoleID.String()))
				return nil, fmt.Errorf("role not found: %w", dto.ErrNotFound)
			}
			user.Role = *role
		}
		user.RoleID = req.RoleID
	}

	if req.ImageURL != nil {
		user.ImageURL = *req.ImageURL
	}
	user.Name = req.Name

	err = us.userRepo.UpdateUserByUserID(ctx, nil, user)
	if err != nil {
		us.logger.Error("failed to update user", zap.String("id", req.ID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to update user: %w", dto.ErrInternal)
	}

	return mapToUserResponse(user, &user.Role), nil
}

func (us *userService) DeleteUserByUserID(ctx context.Context, userID *uuid.UUID) error {
	_, found, err := us.userRepo.GetUserByUserID(ctx, nil, userID)
	if err != nil {
		us.logger.Error("failed to get user by ID", zap.String("userID", userID.String()), zap.Error(err))
		return fmt.Errorf("failed to get user ID: %w", dto.ErrInternal)
	}
	if !found {
		us.logger.Warn("user not found", zap.String("userID", userID.String()))
		return fmt.Errorf("user not found: %w", dto.ErrNotFound)
	}

	if err := us.userRepo.DeleteUserByUserID(ctx, nil, userID); err != nil {
		us.logger.Error("failed to delete user by id", zap.String("userID", userID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete user by id: %w", dto.ErrInternal)
	}

	return nil
}
