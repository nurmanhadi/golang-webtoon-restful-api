package service

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"strconv"
	"time"
	"webtoon/internal/domain/dto"
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"
	"webtoon/internal/infrastructure/storage/s3"
	"webtoon/pkg"
	"webtoon/pkg/image"
	"webtoon/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ComicService interface {
	AddComic(cover *multipart.FileHeader, request *dto.ComicAddRequest) error
	UpdateComic(id string, cover *multipart.FileHeader, request *dto.ComicUpdateRequest) error
	GetById(id string) (*dto.ComicResponse, error)
	GetAll(page string, size string) (*pkg.Paging[[]dto.ComicResponse], error)
	Remove(id string) error
	Search(keyword string, page string, size string) (*pkg.Paging[[]dto.ComicResponse], error)
}

type comicService struct {
	logger          *logrus.Logger
	validation      *validator.Validate
	comicRepository repository.ComicRepository
	s3              s3.S3Storage
}

func NewComicService(logger *logrus.Logger, validation *validator.Validate, comicRepository repository.ComicRepository, s3 s3.S3Storage) ComicService {
	return &comicService{
		logger:          logger,
		validation:      validation,
		comicRepository: comicRepository,
		s3:              s3,
	}
}
func (s *comicService) AddComic(cover *multipart.FileHeader, request *dto.ComicAddRequest) error {
	if err := s.validation.Struct(request); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return err
	}
	if err := image.Validate(cover.Filename); err != nil {
		s.logger.WithError(err).Warn("validation error")
		return response.Exception(400, err.Error())
	}
	filename := fmt.Sprintf("%d.webp", time.Now().Unix())
	go func(filename string) {
		webpFile, err := image.CompressToCwebp(cover)
		if err != nil {
			s.logger.WithError(err).Error("compress avatar error")
			return
		}
		defer webpFile.Close()
		defer os.Remove(webpFile.Name())

		if err := s.s3.UploadFile(webpFile, filename); err != nil {
			s.logger.WithError(err).Error("s3 upload file error")
			return
		}
	}(filename)
	coverUrl := pkg.GenerateUrl(filename)
	id := uuid.NewString()
	comic := &entity.Comic{
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
func (s *comicService) UpdateComic(id string, cover *multipart.FileHeader, request *dto.ComicUpdateRequest) error {
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
		go func(filename string) {
			webpFile, err := image.CompressToCwebp(cover)
			if err != nil {
				s.logger.WithError(err).Error("compress avatar error")
				return
			}
			defer webpFile.Close()
			defer os.Remove(webpFile.Name())

			if err := s.s3.UploadFile(webpFile, comic.CoverFilename); err != nil {
				s.logger.WithError(err).Error("s3 upload file error")
				return
			}
		}(comic.CoverFilename)
	}
	comic.UpdatedAt = time.Now()
	if err := s.comicRepository.Save(comic); err != nil {
		s.logger.WithError(err).Error("comic save error")
	}
	s.logger.WithField("data", id).Info("comic update success")
	return nil
}
func (s *comicService) GetById(id string) (*dto.ComicResponse, error) {
	comic, err := s.comicRepository.FindById(id)
	if err != nil {
		s.logger.WithField("error", id).Warn("comic not found")
		return nil, response.Exception(404, "comic not found")
	}
	genres := make([]dto.GenreResponse, 0, len(comic.ComicGenre))
	for _, comicGenre := range comic.ComicGenre {
		genres = append(genres, dto.GenreResponse{
			Id:   comicGenre.Genre.Id,
			Name: comicGenre.Genre.Name,
		})
	}
	result := &dto.ComicResponse{
		Id:            comic.Id,
		Title:         comic.Title,
		Synopsis:      comic.Synopsis,
		Author:        comic.Author,
		Artist:        comic.Artist,
		Type:          comic.Type,
		CoverFilename: comic.CoverFilename,
		CoverUrl:      comic.CoverUrl,
		CreatedAt:     comic.CreatedAt,
		UpdatedAt:     comic.UpdatedAt,
		Genres:        &genres,
	}
	return result, nil
}
func (s *comicService) GetAll(page string, size string) (*pkg.Paging[[]dto.ComicResponse], error) {
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
	contents := make([]dto.ComicResponse, 0, len(comics))
	for _, comic := range comics {
		contents = append(contents, dto.ComicResponse{
			Id:            comic.Id,
			Title:         comic.Title,
			Synopsis:      comic.Synopsis,
			Author:        comic.Author,
			Artist:        comic.Artist,
			Type:          comic.Type,
			CoverFilename: comic.CoverFilename,
			CoverUrl:      comic.CoverUrl,
			CreatedAt:     comic.CreatedAt,
			UpdatedAt:     comic.UpdatedAt,
		})
	}
	totalComic, err := s.comicRepository.CountTotal()
	if err != nil {
		s.logger.WithError(err).Error("count total comic error")
		return nil, err
	}
	totalPage := int(math.Ceil(float64(totalComic) / float64(newSize)))
	result := &pkg.Paging[[]dto.ComicResponse]{
		Contents:     contents,
		Page:         newPage,
		Size:         newSize,
		TotalPage:    totalPage,
		TotalElement: int(totalComic),
	}
	s.logger.Info("get all comic success")
	return result, nil
}

func (s *comicService) Remove(id string) error {
	comic, err := s.comicRepository.FindById(id)
	if err != nil {
		s.logger.WithField("error", id).Error("comic not found")
		return response.Exception(404, "comic not found")
	}
	go func(filename string) {
		if err := s.s3.RemoveFile(filename); err != nil {
			s.logger.WithError(err).Error("remove file error")
			return
		}
	}(comic.CoverFilename)
	if err := s.comicRepository.Delete(id); err != nil {
		s.logger.WithError(err).Error("comic delete error")
		return err
	}
	s.logger.WithField("data", id).Info("comic delete success")
	return nil
}

func (s *comicService) Search(keyword string, page string, size string) (*pkg.Paging[[]dto.ComicResponse], error) {
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
	comics, err := s.comicRepository.Search(keyword, newPage, newSize)
	if err != nil {
		s.logger.WithError(err).Error("search comic error")
		return nil, err
	}
	totalElement, err := s.comicRepository.CountTotalByKeyword(keyword)
	if err != nil {
		s.logger.WithError(err).Error("count comic by keyword error")
		return nil, err
	}
	contents := make([]dto.ComicResponse, 0, len(comics))
	for _, comic := range comics {
		contents = append(contents, dto.ComicResponse{
			Id:            comic.Id,
			Title:         comic.Title,
			Synopsis:      comic.Synopsis,
			Author:        comic.Author,
			Artist:        comic.Artist,
			Type:          comic.Type,
			CoverFilename: comic.CoverFilename,
			CoverUrl:      comic.CoverUrl,
			CreatedAt:     comic.CreatedAt,
			UpdatedAt:     comic.UpdatedAt,
		})
	}
	totalPage := int(math.Ceil(float64(totalElement) / float64(newSize)))
	result := &pkg.Paging[[]dto.ComicResponse]{
		Contents:     contents,
		Page:         newPage,
		Size:         newSize,
		TotalPage:    totalPage,
		TotalElement: int(totalElement),
	}
	s.logger.WithField("data", keyword).Info("comic search success")
	return result, nil
}
