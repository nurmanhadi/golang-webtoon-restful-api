package service

import (
	"webtoon/internal/domain/dto"
	"webtoon/internal/domain/repository"

	"github.com/sirupsen/logrus"
)

type DashboardService interface {
	Summary() (*dto.DashboardResponse, error)
}
type dashboardService struct {
	logger             *logrus.Logger
	userRepository     repository.UserRepository
	comicRepository    repository.ComicRepository
	chapterRespository repository.ChapterRepository
}

func NewDashboardService(
	logger *logrus.Logger,
	userRepository repository.UserRepository,
	comicRepository repository.ComicRepository,
	chapterRespository repository.ChapterRepository,
) DashboardService {
	return &dashboardService{
		logger:             logger,
		userRepository:     userRepository,
		comicRepository:    comicRepository,
		chapterRespository: chapterRespository,
	}
}
func (s *dashboardService) Summary() (*dto.DashboardResponse, error) {
	totalUser, err := s.userRepository.CountTotalUser()
	if err != nil {
		s.logger.WithError(err).Error("count total user error")
		return nil, err
	}
	totalComic, err := s.comicRepository.CountTotal()
	if err != nil {
		s.logger.WithError(err).Error("count total comic error")
		return nil, err
	}
	totalView, err := s.comicRepository.CountTotalView()
	if err != nil {
		s.logger.WithError(err).Error("count total view error")
		return nil, err
	}
	totalViewDaily, err := s.comicRepository.CountTotalViewDaily()
	if err != nil {
		s.logger.WithError(err).Error("count total view daily error")
		return nil, err
	}
	totalViewWeekly, err := s.comicRepository.CountTotalViewWeekly()
	if err != nil {
		s.logger.WithError(err).Error("count total view weekly error")
		return nil, err
	}
	totalViewMonthly, err := s.comicRepository.CountTotalViewMonthly()
	if err != nil {
		s.logger.WithError(err).Error("count total view monthly error")
		return nil, err
	}
	totalChapter, err := s.chapterRespository.CountTotalChapter()
	if err != nil {
		s.logger.WithError(err).Error("count total chapter error")
		return nil, err
	}
	result := &dto.DashboardResponse{
		TotalComic:        int(totalComic),
		TotalChapter:      int(totalChapter),
		TotalUser:         int(totalUser),
		TotalViews:        int(totalView),
		TotalViewsDaily:   int(totalViewDaily),
		TotalViewsWeekly:  int(totalViewWeekly),
		TotalViewsMonthly: int(totalViewMonthly),
	}
	s.logger.Info("get summary success")
	return result, nil
}
