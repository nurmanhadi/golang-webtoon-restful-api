package repository

import "webtoon/internal/domain/entity"

type ComicRepository interface {
	Save(comic *entity.Comic) error
	FindById(id string) (*entity.Comic, error)
	FindAll(page int, size int) ([]entity.Comic, error)
	CountTotal() (int64, error)
	Delete(id string) error
	Search(key string, page int, size int) ([]entity.Comic, error)
	CountTotalByKeyword(key string) (int64, error)
	CountById(id string) (int64, error)
	FindAllByType(comicType string, page int, size int) ([]entity.Comic, error)
	CountTotalByType(comicType string) (int64, error)
	FindAllByViewsByPeriod(limit int, timePeriod string) ([]entity.Comic, error)
}
