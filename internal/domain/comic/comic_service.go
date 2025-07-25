package comic

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"strconv"
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
	UpdateComic(id string, cover *multipart.FileHeader, request *ComicUpdateRequest) error
	GetById(id string) (*ComicResponse, error)
	GetAll(page string, size string) (*pkg.Paging[[]ComicResponse], error)
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
	s.logger.WithField("data", id).Info("comic save success")
	return nil
}
func (s *comicService) UpdateComic(id string, cover *multipart.FileHeader, request *ComicUpdateRequest) error {
	if err := s.validation.Struct(request); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return err
	}
	comic, err := s.comicRepository.FindById(id)
	if err != nil {
		s.logger.WithField("error", id).Warn("comic not found")
		return response.Exception(404, "comic not found")
	}
	if request.Title != "" {
		comic.Title = request.Title
	}
	if request.Synopsis != "" {
		comic.Synopsis = request.Synopsis
	}
	if request.Author != "" {
		comic.Author = request.Author
	}
	if request.Artist != "" {
		comic.Artist = request.Artist
	}
	if request.Type != "" {
		comic.Type = request.Type
	}
	if cover != nil {
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

		if err := s.s3.UploadFile(webpFile, comic.CoverFilename); err != nil {
			s.logger.WithError(err).Error("s3 upload file error")
			return err
		}
	}
	comic.UpdatedAt = time.Now()
	if err := s.comicRepository.Save(comic); err != nil {
		s.logger.WithError(err).Error("comic save error")
	}
	s.logger.WithField("data", id).Info("comic update success")
	return nil
}
func (s *comicService) GetById(id string) (*ComicResponse, error) {
	comic, err := s.comicRepository.FindById(id)
	if err != nil {
		s.logger.WithField("error", id).Warn("comic not found")
		return nil, response.Exception(404, "comic not found")
	}
	result := ComicResponse(*comic)
	return &result, nil
}
func (s *comicService) GetAll(page string, size string) (*pkg.Paging[[]ComicResponse], error) {
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
	comics, err := s.comicRepository.FindAll(newPage, newSize)
	if err != nil {
		s.logger.WithError(err).Error("find all comic error")
		return nil, err
	}
	contents := make([]ComicResponse, 0, len(comics))
	for _, comic := range comics {
		contents = append(contents, ComicResponse(comic))
	}
	totalComic, err := s.comicRepository.CountTotal()
	if err != nil {
		s.logger.WithError(err).Error("count total comic error")
		return nil, err
	}
	totalPage := int(math.Ceil(float64(totalComic) / float64(newSize)))
	result := &pkg.Paging[[]ComicResponse]{
		Contents:     contents,
		Page:         newPage,
		Size:         newSize,
		TotalPage:    totalPage,
		TotalElement: int(totalComic),
	}
	return result, nil
}
