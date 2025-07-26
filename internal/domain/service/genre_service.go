package service

import (
	"strconv"
	"webtoon/internal/domain/dto"
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"
	"webtoon/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type GenreService interface {
	AddGenre(request *dto.GenreAddRequest) error
	GetAll() ([]dto.GenreResponse, error)
	Remove(id string) error
}
type genreService struct {
	logger          *logrus.Logger
	validation      *validator.Validate
	genreRepository repository.GenreRepository
}

func NewGenreService(logger *logrus.Logger, validation *validator.Validate, genreRepository repository.GenreRepository) GenreService {
	return &genreService{
		logger:          logger,
		validation:      validation,
		genreRepository: genreRepository,
	}
}

func (s *genreService) AddGenre(request *dto.GenreAddRequest) error {
	if err := s.validation.Struct(request); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return err
	}
	genre := &entity.Genre{
		Name: request.Name,
	}
	if err := s.genreRepository.Save(genre); err != nil {
		s.logger.WithError(err).Error("genre save error")
		return err
	}
	s.logger.Info("add genre success")
	return nil
}
func (s *genreService) GetAll() ([]dto.GenreResponse, error) {
	genres, err := s.genreRepository.FindAll()
	if err != nil {
		s.logger.WithError(err).Error("find all genres error")
		return nil, err
	}
	var result []dto.GenreResponse
	for _, genre := range genres {
		result = append(result, dto.GenreResponse{
			Id:   genre.Id,
			Name: genre.Name,
		})
	}
	s.logger.Info("get all genres success")
	return result, nil
}
func (s *genreService) Remove(id string) error {
	newId, err := strconv.Atoi(id)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return response.Exception(400, "id most be number")
	}
	count, err := s.genreRepository.Count(newId)
	if err != nil {
		s.logger.WithError(err).Warn("count genre error")
		return err
	}
	if count < 1 {
		s.logger.WithField("error", id).Warn("genre not found")
		return response.Exception(404, "genre not found")
	}
	if err := s.genreRepository.Remove(newId); err != nil {
		s.logger.WithError(err).Error("genre remove error")
		return err
	}
	s.logger.Info("remove genres success")
	return nil
}
