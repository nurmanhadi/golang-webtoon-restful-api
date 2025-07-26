package service

import (
	"math"
	"strconv"
	"webtoon/internal/domain/dto"
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"
	"webtoon/pkg"
	"webtoon/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type GenreService interface {
	AddGenre(request *dto.GenreAddRequest) error
	GetAll() ([]dto.GenreResponse, error)
	Remove(id string) error
	GetById(id string, page string, size string) (*dto.GenreResponse, error)
}
type genreService struct {
	logger               *logrus.Logger
	validation           *validator.Validate
	genreRepository      repository.GenreRepository
	comicGenreRepository repository.ComicGenreRepository
}

func NewGenreService(logger *logrus.Logger, validation *validator.Validate, genreRepository repository.GenreRepository, comicGenreRepository repository.ComicGenreRepository) GenreService {
	return &genreService{
		logger:               logger,
		validation:           validation,
		genreRepository:      genreRepository,
		comicGenreRepository: comicGenreRepository,
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
func (s *genreService) GetById(id string, page string, size string) (*dto.GenreResponse, error) {
	newId, err := strconv.Atoi(id)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return nil, response.Exception(400, "id most be number")
	}
	newPage, err := strconv.Atoi(page)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return nil, response.Exception(400, "page most be number")
	}
	newSize, err := strconv.Atoi(size)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return nil, response.Exception(400, "size most be number")
	}
	genre, err := s.genreRepository.FindById(newId)
	if err != nil {
		s.logger.WithField("error", id).Warn("genre not found")
		return nil, response.Exception(404, "genre not found")
	}
	comicGenres, err := s.comicGenreRepository.FindAllByGenreId(newId, newPage, newSize)
	if err != nil {
		s.logger.WithError(err).Error("find all comic genres error")
		return nil, err
	}
	totalElement, err := s.comicGenreRepository.CountByGenreId(newId)
	if err != nil {
		s.logger.WithError(err).Error("count comic genres error")
		return nil, err
	}
	comics := make([]dto.ComicResponse, 0, len(comicGenres))
	for _, comicGenre := range comicGenres {
		comics = append(comics, dto.ComicResponse{
			Id:            comicGenre.Comic.Id,
			Title:         comicGenre.Comic.Title,
			Synopsis:      comicGenre.Comic.Synopsis,
			Author:        comicGenre.Comic.Author,
			Artist:        comicGenre.Comic.Artist,
			Type:          comicGenre.Comic.Type,
			CoverFilename: comicGenre.Comic.CoverFilename,
			CoverUrl:      comicGenre.Comic.CoverUrl,
			CreatedAt:     comicGenre.Comic.CreatedAt,
			UpdatedAt:     comicGenre.Comic.UpdatedAt,
		})
	}
	totalPage := int(math.Ceil(float64(totalElement) / float64(newSize)))
	result := &dto.GenreResponse{
		Id:   genre.Id,
		Name: genre.Name,
		Comics: &pkg.Paging[[]dto.ComicResponse]{
			Contents:     comics,
			Page:         newPage,
			Size:         newSize,
			TotalPage:    totalPage,
			TotalElement: int(totalElement),
		},
	}
	s.logger.WithField("data", id).Info("get genre success")
	return result, nil
}
func (s *genreService) Remove(id string) error {
	newId, err := strconv.Atoi(id)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return response.Exception(400, "id most be number")
	}
	countGenre, err := s.genreRepository.Count(newId)
	if err != nil {
		s.logger.WithError(err).Warn("count genre error")
		return err
	}
	if countGenre < 1 {
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
