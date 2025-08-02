package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"time"
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"
	"webtoon/internal/infrastructure/storage/s3"
	"webtoon/pkg"
	"webtoon/pkg/image"
	"webtoon/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ContentService interface {
	AddBulkContent(chapterId string, contents []*multipart.FileHeader) error
	Remove(id string) error
}
type contentService struct {
	logger            *logrus.Logger
	validation        *validator.Validate
	contentRepository repository.ContentRepository
	chapterRepository repository.ChapterRepository
	s3                s3.S3Storage
}

func NewContentService(
	logger *logrus.Logger,
	validation *validator.Validate,
	contentRepository repository.ContentRepository,
	chapterRepository repository.ChapterRepository,
	s3 s3.S3Storage,
) ContentService {
	return &contentService{
		logger:            logger,
		validation:        validation,
		contentRepository: contentRepository,
		chapterRepository: chapterRepository,
		s3:                s3,
	}
}
func (s *contentService) AddBulkContent(chapterId string, contents []*multipart.FileHeader) error {
	newChapterId, err := strconv.Atoi(chapterId)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return response.Exception(400, "chapter id most be number")
	}
	countChapter, err := s.chapterRepository.Count(newChapterId)
	if err != nil {
		s.logger.WithError(err).Warn("count chapter error")
		return err
	}
	if countChapter < 1 {
		s.logger.WithField("error", chapterId).Warn("chapter not found")
		return response.Exception(404, "chapter not found")
	}

	for _, content := range contents {
		if err := image.Validate(content.Filename); err != nil {
			s.logger.WithError(err).Warn("validation error")
			return response.Exception(400, err.Error())
		}
	}

	count := 0
	for _, content := range contents {
		filename := fmt.Sprintf("%d.webp", time.Now().UnixNano())
		webpFile, err := image.CompressToCwebp(content)
		if err != nil {
			s.logger.WithError(err).Error("compress avatar error")
			return err
		}
		defer webpFile.Close()
		defer os.Remove(webpFile.Name())

		if err := s.s3.UploadFile(webpFile, filename); err != nil {
			s.logger.WithError(err).Error("s3 upload file error")
			return err
		}
		count += 1
		fmt.Println("upload: ", +count)

		url := pkg.GenerateUrl(filename)
		contentEnt := &entity.Content{
			ChapterId: int64(newChapterId),
			Filename:  filename,
			Url:       url,
		}
		if err := s.contentRepository.Save(contentEnt); err != nil {
			s.logger.WithError(err).Error("genre save error")
			return err
		}
	}
	s.logger.Info("add bulk content success")
	return nil
}
func (s *contentService) Remove(id string) error {
	newId, err := strconv.Atoi(id)
	if err != nil {
		s.logger.WithError(err).Warn("parse string to int error")
		return response.Exception(400, "id most be number")
	}
	content, err := s.contentRepository.FindById(newId)
	if err != nil {
		s.logger.WithField("error", id).Warn("content not found")
		return response.Exception(404, "content not found")
	}
	go func(filename string) {
		if err := s.s3.RemoveFile(filename); err != nil {
			s.logger.WithError(err).Error("s3 remove file error")
			return
		}
	}(content.Filename)
	if err := s.contentRepository.Delete(newId); err != nil {
		s.logger.WithError(err).Error("content remove error")
		return err
	}
	s.logger.WithField("data", id).Info("content remove success")
	return nil
}
