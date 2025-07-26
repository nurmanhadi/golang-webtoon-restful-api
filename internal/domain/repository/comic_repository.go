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
}
