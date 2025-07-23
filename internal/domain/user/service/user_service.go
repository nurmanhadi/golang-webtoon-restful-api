package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"webtoon/internal/domain/user/dto"
	"webtoon/internal/domain/user/repository"
	"webtoon/pkg/image"
	"webtoon/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	GetById(id string) (*dto.UserResponse, error)
	UpdateUsername(id string, request dto.UserUpdateUsernameRequest) error
	UploadAvatar(id string, avatar multipart.FileHeader) error
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
func (s *service) UpdateUsername(id string, request dto.UserUpdateUsernameRequest) error {
	if err := s.validation.Struct(&request); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return err
	}
	user, err := s.userRepository.FindById(id)
	if err != nil {
		s.logger.WithField("error", id).Warn("user not found")
		return response.Exception(404, "user not found")
	}
	if user.Username == request.Username {
		s.logger.WithField("data", id).Info("update username success")
		return nil
	}
	user.Username = request.Username
	if err := s.userRepository.Save(user); err != nil {
		s.logger.WithError(err).Error("save user to database error")
		return err
	}
	s.logger.WithField("data", id).Info("update username success")
	return nil
}
func (s *service) UploadAvatar(id string, avatar multipart.FileHeader) error {
	if err := image.Validate(avatar); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return response.Exception(404, err.Error())
	}
	webpFile, err := image.CompressToCwebp(avatar)
	if err != nil {
		s.logger.WithError(err).Error("compress avatar error")
		return err
	}
	defer webpFile.Close()
	os.Remove(webpFile.Name())
	fmt.Println(webpFile.Name())

	return nil
}
