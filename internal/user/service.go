package user

import (
	"context"
	"fmt"
	"rest-api/pkg/hasher"
	"rest-api/pkg/logging"
)

type Service interface {
	Create(ctx context.Context, dto CreateUserDTO) (User, error)
	GetAllUsers(ctx context.Context) (users []User, err error)
	GetUser(ctx context.Context, id string) (User, error)
	UpdateUser(ctx context.Context, dto UpdateUserDTO) error
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	storage Storage
	logger  *logging.Logger
}

func NewUserService(storage Storage, logger *logging.Logger) *userService {
	return &userService{
		storage: storage,
		logger:  logger,
	}
}

func (s *userService) Create(ctx context.Context, dto CreateUserDTO) (User, error) {
	var user User
	passwordHash, err := hasher.Encrypt(dto.Password)
	if err != nil {
		return user, fmt.Errorf("Encrypting body error: %v", err)
	}

	user.Username = dto.Username
	user.PasswordHash = passwordHash
	user.Email = dto.Email

	id, err := s.storage.Create(ctx, user)
	if err != nil {
		return user, err
	}
	user.ID = id

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (User, error) {
	user, err := s.storage.FindOne(ctx, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) (users []User, err error) {
	users, err = s.storage.FindAll(ctx)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (s *userService) UpdateUser(ctx context.Context, dto UpdateUserDTO) error {
	passwordHash, err := hasher.Encrypt(dto.Password)
	if err != nil {
		return fmt.Errorf("Encrypting body error: %v", err)
	}

	err = s.storage.Update(ctx, User{
		ID:           dto.ID,
		Username:     dto.Username,
		PasswordHash: passwordHash,
		Email:        dto.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	err := s.storage.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
