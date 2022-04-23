package user

import (
	"rest-api/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

// func (s *Service) CreateUser(ctx context.Context, dto CreateUserDTO) (User, error) {

// 	id, err := s.storage.Create(ctx, User{})
// 	if err != nil {
// 		return nil, err
// 	}
// }
