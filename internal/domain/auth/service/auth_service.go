package service

import (
	"webtoon/internal/domain/auth/dto"
	"webtoon/internal/domain/auth/repository"
	"webtoon/internal/domain/user/entity"
	"webtoon/pkg/response"
	"webtoon/pkg/role"
	"webtoon/pkg/security"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(request dto.AuthRequest) error
	Login(request dto.AuthRequest) (*dto.AuthResponse, error)
}
type service struct {
	logger         *logrus.Logger
	validation     *validator.Validate
	authRepository repository.AuthRepository
}

func NewAuthService(logger *logrus.Logger, validation *validator.Validate, authRepository repository.AuthRepository) AuthService {
	return &service{
		logger:         logger,
		validation:     validation,
		authRepository: authRepository,
	}
}
func (s *service) Register(request dto.AuthRequest) error {
	if err := s.validation.Struct(&request); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return err
	}
	userCount, err := s.authRepository.CountByUsername(request.Username)
	if err != nil {
		s.logger.WithError(err).Error("count user error")
		return err
	}
	if userCount > 0 {
		s.logger.WithField("error", userCount).Warn("username already exists")
		return response.Exception(400, "username already exists")
	}
	newPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("hash password error")
		return err
	}
	user := &entity.User{
		Id:       uuid.NewString(),
		Username: request.Username,
		Password: string(newPassword),
		Role:     string(role.USER),
	}

	if err := s.authRepository.Save(user); err != nil {
		s.logger.WithError(err).Error("save user to database error")
		return err
	}
	s.logger.WithField("data", user.Username).Info("registration success")
	return nil
}

func (s *service) Login(request dto.AuthRequest) (*dto.AuthResponse, error) {
	if err := s.validation.Struct(&request); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return nil, err
	}
	user, err := s.authRepository.FindByUsername(request.Username)
	if err != nil {
		s.logger.WithField("error", request.Username).Warn("username or password wrong")
		return nil, response.Exception(400, "username or password wrong")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		s.logger.WithField("error", request.Username).Warn("username or password wrong")
		return nil, response.Exception(400, "username or password wrong")
	}
	token, err := security.JwtGenerateAccessToken(user.Id, user.Role)
	if err != nil {
		s.logger.WithError(err).Error("generate access token error")
		return nil, err
	}
	response := &dto.AuthResponse{
		AccessToken: token,
	}
	s.logger.WithField("data", user.Username).Info("login success")
	return response, nil
}
