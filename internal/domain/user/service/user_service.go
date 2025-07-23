package service

import (
	"webtoon/internal/domain/user/dto"
	"webtoon/internal/domain/user/repository"
	"webtoon/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	GetById(id string) (*dto.UserResponse, error)
}
type service struct {
	logger         *logrus.Logger
	validation     *validator.Validate
	userRepository repository.UserRepository
}

func NewUserService(logger *logrus.Logger, validation *validator.Validate, userRepository repository.UserRepository) UserService {
	return &service{
		logger:         logger,
		validation:     validation,
		userRepository: userRepository,
	}
}
func (s *service) GetById(id string) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindById(id)
	if err != nil {
		s.logger.WithField("error", id).Warn("user not found")
		return nil, response.Exception(404, "user not found")
	}
	response := &dto.UserResponse{
		Id:             user.Id,
		Username:       user.Username,
		AvatarFilename: user.AvatarFilename,
		AvatarUrl:      user.AvatarUrl,
	}
	s.logger.WithField("data", user.Id).Info("get user by id success")
	return response, nil
}
