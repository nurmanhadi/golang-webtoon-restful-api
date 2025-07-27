package repository

import "webtoon/internal/domain/entity"

type ChapterRepository interface {
	Save(chapter *entity.Chapter) error
	FindById(id int) (*entity.Chapter, error)
	FindByIdAndNumber(id int, number int) (*entity.Chapter, error)
}
