package user

import (
	"fmt"
	"mime/multipart"
	"os"
	"time"
	"webtoon/internal/infrastructure/storage/s3"
	"webtoon/pkg"
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
	s3             s3.S3Storage
}

func NewUserService(logger *logrus.Logger, validation *validator.Validate, userRepository UserRepository, s3 s3.S3Storage) UserService {
	return &userService{
		logger:         logger,
		validation:     validation,
		userRepository: userRepository,
		s3:             s3,
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
	user, err := s.userRepository.FindById(id)
	if err != nil {
		s.logger.WithField("error", id).Warn("user not found")
		return response.Exception(404, "user not found")
	}
	if err := image.Validate(avatar.Filename); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return response.Exception(400, err.Error())
	}
	webpFile, err := image.CompressToCwebp(avatar)
	if err != nil {
		s.logger.WithError(err).Error("compress avatar error")
		return err
	}
	defer webpFile.Close()
	defer os.Remove(webpFile.Name())
	if user.AvatarFilename == "" && user.AvatarUrl == "" {
		filename := fmt.Sprintf("%d.webp", time.Now().Unix())
		if err := s.s3.UploadFile(webpFile, filename); err != nil {
			s.logger.WithError(err).Error("s3 upload file error")
			return err
		}
		avatarUrl := pkg.GenerateUrl(filename)
		user.AvatarFilename = filename
		user.AvatarUrl = avatarUrl
		if err := s.userRepository.Save(user); err != nil {
			s.logger.WithError(err).Error("save user to database error")
			return err
		}
	} else {
		if err := s.s3.UploadFile(webpFile, user.AvatarFilename); err != nil {
			s.logger.WithError(err).Error("s3 upload file error")
			return err
		}
	}
	s.logger.WithField("data", id).Info("upload avatar success")
	return nil
}
