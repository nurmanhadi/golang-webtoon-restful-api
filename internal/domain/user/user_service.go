package user

import (
	"fmt"
	"mime/multipart"
	"os"
	"webtoon/pkg/image"
	"webtoon/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	GetById(id string) (*UserResponse, error)
	UpdateUsername(id string, request UserUpdateUsernameRequest) error
	UploadAvatar(id string, avatar multipart.FileHeader) error
}
type userService struct {
	logger         *logrus.Logger
	validation     *validator.Validate
	userRepository UserRepository
}

func NewUserService(logger *logrus.Logger, validation *validator.Validate, userRepository UserRepository) UserService {
	return &userService{
		logger:         logger,
		validation:     validation,
		userRepository: userRepository,
	}
}
func (s *userService) GetById(id string) (*UserResponse, error) {
	user, err := s.userRepository.FindById(id)
	if err != nil {
		s.logger.WithField("error", id).Warn("user not found")
		return nil, response.Exception(404, "user not found")
	}
	response := &UserResponse{
		Id:             user.Id,
		Username:       user.Username,
		AvatarFilename: user.AvatarFilename,
		AvatarUrl:      user.AvatarUrl,
	}
	s.logger.WithField("data", user.Id).Info("get user by id success")
	return response, nil
}
func (s *userService) UpdateUsername(id string, request UserUpdateUsernameRequest) error {
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
func (s *userService) UploadAvatar(id string, avatar multipart.FileHeader) error {
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
