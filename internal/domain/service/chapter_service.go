package service

import (
	"sort"
	"strconv"
	"time"
	"webtoon/internal/domain/dto"
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"
	"webtoon/internal/infrastructure/storage/s3"
	"webtoon/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ChapterService interface {
	AddChapter(request *dto.ChapterAddRequest) error
	UpdateChapter(id string, request dto.ChapterUpdateRequest) error
	GetByComicIdAndNumber(comicId string, number string) (*dto.ChapterResponse, error)
	Remove(id string) error
}
type chapterService struct {
	logger            *logrus.Logger
	validation        *validator.Validate
	chapterRepository repository.ChapterRepository
	comicRepository   repository.ComicRepository
	s3                s3.S3Storage
}

func NewChapterService(logger *logrus.Logger, validation *validator.Validate, chapterRepository repository.ChapterRepository, comicRepository repository.ComicRepository, s3 s3.S3Storage) ChapterService {
	return &chapterService{
		logger:            logger,
		validation:        validation,
		chapterRepository: chapterRepository,
		comicRepository:   comicRepository,
		s3:                s3,
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
	comic, err := s.comicRepository.FindById(chapter.ComicId)
	if err != nil {
		s.logger.WithField("error", id).Warn("comic not found")
		return response.Exception(404, "comic not found")
	}
	if request.Number != nil {
		chapter.Number = *request.Number
	}
	if request.Publish != nil {
		chapter.Publish = *request.Publish
		comic.UpdatedPost = time.Now()
	}
	if err := s.chapterRepository.Save(chapter); err != nil {
		s.logger.WithError(err).Error("chapter save error")
		return err
	}
	if err := s.comicRepository.Save(comic); err != nil {
		s.logger.WithError(err).Error("comic save error")
		return err
	}
	s.logger.WithField("data", id).Info("update chapter success")
	return nil
}
func (s *chapterService) GetByComicIdAndNumber(comicId string, number string) (*dto.ChapterResponse, error) {
	countComic, err := s.comicRepository.CountById(comicId)
	if err != nil {
		s.logger.WithError(err).Warn("count comic error")
		return nil, err
	}
	if countComic < 1 {
		s.logger.WithField("error", comicId).Warn("comic not found")
		return nil, response.Exception(404, "comic not found")
	}
	newNumber, err := strconv.Atoi(number)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return nil, response.Exception(400, "number most be number")
	}

	chapter, err := s.chapterRepository.FindByComicIdAndNumber(comicId, newNumber)
	if err != nil {
		s.logger.WithField("error", number).Warn("chapter not found")
		return nil, response.Exception(404, "chapter not found")
	}

	chapters := make([]dto.ChapterResponse, 0, len(chapter.Comic.Chapters))
	if len(chapter.Comic.Chapters) != 0 {
		for _, chapter := range chapter.Comic.Chapters {
			chapters = append(chapters, dto.ChapterResponse{
				Id:        chapter.Id,
				ComicId:   chapter.ComicId,
				Number:    chapter.Number,
				Publish:   chapter.Publish,
				CreatedAt: chapter.CreatedAt,
			})
		}
		sort.Slice(chapters, func(i, j int) bool {
			return chapters[i].Number > chapters[j].Number // descending
		})
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
		Chapters:      &chapters,
	}

	contents := make([]dto.ContentResponse, 0, len(chapter.Contents))
	if len(chapter.Contents) != 0 {
		for _, content := range chapter.Contents {
			contents = append(contents, dto.ContentResponse{
				Id:        content.Id,
				ChapterId: content.ChapterId,
				Filename:  content.Filename,
				Url:       content.Url,
			})
		}
		sort.Slice(contents, func(i, j int) bool {
			return contents[i].Filename < contents[j].Filename // ascending
		})
	}
	result := &dto.ChapterResponse{
		Id:        chapter.Id,
		ComicId:   chapter.ComicId,
		Number:    chapter.Number,
		Publish:   chapter.Publish,
		CreatedAt: chapter.CreatedAt,
		Comic:     comic,
		Contents:  &contents,
	}
	s.logger.WithField("data", number).Info("get by id and number chapter success")
	return result, nil
}
func (s *chapterService) Remove(id string) error {
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
	if len(chapter.Contents) != 0 {
		for _, content := range chapter.Contents {
			if err := s.s3.RemoveFile(content.Filename); err != nil {
				s.logger.WithError(err).Error("s3 remove file error")
				return err
			}
		}
	}
	if err := s.chapterRepository.Delete(newId); err != nil {
		s.logger.WithError(err).Error("chapter remove error")
		return err
	}
	s.logger.WithField("data", id).Info("chapter remove success")
	return nil
}
