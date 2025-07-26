package repository

import "webtoon/internal/domain/entity"

type GenreRepository interface {
	Save(genre *entity.Genre) error
	Remove(id int) error
	FindAll() ([]entity.Genre, error)
	Count(id int) (int64, error)
}
