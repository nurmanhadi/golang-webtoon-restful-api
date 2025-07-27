package service

import (
	"fmt"
	"log"
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
		return response.Exception(400, "id most be number")
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
		log.Println("uploaded")
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
