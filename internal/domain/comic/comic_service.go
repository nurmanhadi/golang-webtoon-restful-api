package comic

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
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ComicService interface {
	AddComic(cover *multipart.FileHeader, request *ComicAddRequest) error
}

type comicService struct {
	logger          *logrus.Logger
	validation      *validator.Validate
	comicRepository ComicRepository
	s3              s3.S3Storage
}

func NewComicService(logger *logrus.Logger, validation *validator.Validate, comicRepository ComicRepository, s3 s3.S3Storage) ComicService {
	return &comicService{
		logger:          logger,
		validation:      validation,
		comicRepository: comicRepository,
		s3:              s3,
	}
}
func (s *comicService) AddComic(cover *multipart.FileHeader, request *ComicAddRequest) error {
	if err := s.validation.Struct(request); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return err
	}
	if err := image.Validate(cover.Filename); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return response.Exception(400, err.Error())
	}
	webpFile, err := image.CompressToCwebp(cover)
	if err != nil {
		s.logger.WithError(err).Error("compress avatar error")
		return err
	}
	defer webpFile.Close()
	defer os.Remove(webpFile.Name())

	filename := fmt.Sprintf("%d.webp", time.Now().Unix())
	if err := s.s3.UploadFile(webpFile, filename); err != nil {
		s.logger.WithError(err).Error("s3 upload file error")
		return err
	}
	coverUrl := pkg.GenerateUrl(filename)
	id := uuid.NewString()
	comic := &Comic{
		Id:            id,
		Title:         request.Title,
		Synopsis:      request.Synopsis,
		Author:        request.Author,
		Artist:        request.Artist,
		Type:          request.Type,
		CoverFilename: filename,
		CoverUrl:      coverUrl,
	}
	if err := s.comicRepository.Save(comic); err != nil {
		s.logger.WithError(err).Error("comic save error")
	}
	s.logger.WithField("data", id).Error("comic save success")
	return nil
}
