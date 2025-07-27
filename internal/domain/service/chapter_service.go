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

type ChapterService interface {
	AddChapter(request *dto.ChapterAddRequest) error
	UpdateChapter(id string, request dto.ChapterUpdateRequest) error
	GetByIdAndComicId(id string, number string) (*dto.ChapterResponse, error)
}
type chapterService struct {
	logger            *logrus.Logger
	validation        *validator.Validate
	chapterRepository repository.ChapterRepository
	comicRepository   repository.ComicRepository
}

func NewChapterService(logger *logrus.Logger, validation *validator.Validate, chapterRepository repository.ChapterRepository, comicRepository repository.ComicRepository) ChapterService {
	return &chapterService{
		logger:            logger,
		validation:        validation,
		chapterRepository: chapterRepository,
		comicRepository:   comicRepository,
	}
}
func (s *chapterService) AddChapter(request *dto.ChapterAddRequest) error {
	if err := s.validation.Struct(request); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return err
	}
	countComic, err := s.comicRepository.CountById(request.ComicId)
	if err != nil {
		s.logger.WithError(err).Warn("count comic error")
		return err
	}
	if countComic < 1 {
		s.logger.WithField("error", request.ComicId).Warn("comic not found")
		return response.Exception(404, "comic not found")
	}
	chapter := &entity.Chapter{
		ComicId: request.ComicId,
		Number:  request.Number,
		Publish: false,
	}
	if err := s.chapterRepository.Save(chapter); err != nil {
		s.logger.WithError(err).Error("chapter save error")
		return err
	}
	s.logger.Info("add chapter success")
	return nil
}
func (s *chapterService) UpdateChapter(id string, request dto.ChapterUpdateRequest) error {
	if err := s.validation.Struct(request); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return err
	}
	newId, err := strconv.Atoi(id)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return response.Exception(400, "id most be number")
	}
	chapter, err := s.chapterRepository.FindById(newId)
	if err != nil {
		s.logger.WithField("error", id).Warn("chapter not found")
		return response.Exception(404, "chapter not found")
	}
	if request.Number != nil {
		chapter.Number = *request.Number
	}
	if request.Publish != nil {
		chapter.Publish = *request.Publish
	}
	if err := s.chapterRepository.Save(chapter); err != nil {
		s.logger.WithError(err).Error("chapter save error")
		return err
	}
	s.logger.WithField("data", id).Info("update chapter success")
	return nil
}
func (s *chapterService) GetByIdAndComicId(id string, number string) (*dto.ChapterResponse, error) {
	newId, err := strconv.Atoi(id)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return nil, response.Exception(400, "id most be number")
	}
	newNumber, err := strconv.Atoi(number)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return nil, response.Exception(400, "number most be number")
	}
	chapter, err := s.chapterRepository.FindByIdAndNumber(newId, newNumber)
	if err != nil {
		s.logger.WithField("error", id).Warn("chapter not found")
		return nil, response.Exception(404, "chapter not found")
	}
	comic := &dto.ComicResponse{
		Id:            chapter.Comic.Id,
		Title:         chapter.Comic.Title,
		Synopsis:      chapter.Comic.Synopsis,
		Author:        chapter.Comic.Author,
		Artist:        chapter.Comic.Artist,
		Type:          chapter.Comic.Type,
		CoverFilename: chapter.Comic.CoverFilename,
		CoverUrl:      chapter.Comic.CoverUrl,
		CreatedAt:     chapter.Comic.CreatedAt,
		UpdatedAt:     chapter.Comic.UpdatedAt,
	}
	result := &dto.ChapterResponse{
		Id:        chapter.Id,
		ComicId:   chapter.ComicId,
		Number:    chapter.Number,
		Publish:   chapter.Publish,
		CreatedAt: chapter.CreatedAt,
		Comic:     comic,
	}
	s.logger.WithField("data", id).Info("get by id and number chapter success")
	return result, nil
}
