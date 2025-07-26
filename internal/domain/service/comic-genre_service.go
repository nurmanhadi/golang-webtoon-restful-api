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

type ComicGenreService interface {
	AddComicGenre(request *dto.ComicGenreAddRequest) error
	RemoveComicGenreByComicIdAndGenreId(comicId string, genreId string) error
}
type comicGenreService struct {
	logger               *logrus.Logger
	validation           *validator.Validate
	comicGenreRepository repository.ComicGenreRepository
	comicRepository      repository.ComicRepository
	genreRepository      repository.GenreRepository
}

func NewComicGenreService(logger *logrus.Logger, validation *validator.Validate, comicGenreRepository repository.ComicGenreRepository, comicRepository repository.ComicRepository, genreRepository repository.GenreRepository) ComicGenreService {
	return &comicGenreService{
		logger:               logger,
		validation:           validation,
		comicGenreRepository: comicGenreRepository,
		comicRepository:      comicRepository,
		genreRepository:      genreRepository,
	}
}
func (s *comicGenreService) AddComicGenre(request *dto.ComicGenreAddRequest) error {
	if err := s.validation.Struct(request); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return err
	}
	_, err := s.comicRepository.FindById(request.ComicId)
	if err != nil {
		s.logger.WithField("error", request.ComicId).Error("comic not found")
		return response.Exception(404, "comic not found")
	}
	countGenre, err := s.genreRepository.Count(request.GenreId)
	if err != nil {
		s.logger.WithError(err).Warn("count genre error")
		return err
	}
	if countGenre < 1 {
		s.logger.WithField("error", request.GenreId).Warn("genre not found")
		return response.Exception(404, "genre not found")
	}
	comicGenre := &entity.ComicGenre{
		ComicId: request.ComicId,
		GenreId: request.GenreId,
	}
	if err := s.comicGenreRepository.Save(comicGenre); err != nil {
		s.logger.WithError(err).Error("comic genre save error")
		return err
	}
	s.logger.Info("add comic genre success")
	return nil
}
func (s *comicGenreService) RemoveComicGenreByComicIdAndGenreId(comicId string, genreId string) error {
	_, err := s.comicRepository.FindById(comicId)
	if err != nil {
		s.logger.WithField("error", comicId).Error("comic not found")
		return response.Exception(404, "comic not found")
	}
	newGenreId, err := strconv.Atoi(genreId)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return response.Exception(400, "id most be number")
	}
	countGenre, err := s.genreRepository.Count(newGenreId)
	if err != nil {
		s.logger.WithError(err).Warn("count genre error")
		return err
	}
	if countGenre < 1 {
		s.logger.WithField("error", newGenreId).Warn("genre not found")
		return response.Exception(404, "genre not found")
	}
	if err := s.comicGenreRepository.DeleteByComicIdAndGenreId(comicId, newGenreId); err != nil {
		s.logger.WithError(err).Error("comic genre delete error")
		return err
	}
	s.logger.Info("remove comic genre success")
	return nil
}
