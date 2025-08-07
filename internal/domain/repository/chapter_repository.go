package repository

import "webtoon/internal/domain/entity"

type ChapterRepository interface {
	Save(chapter *entity.Chapter) error
	FindById(id int) (*entity.Chapter, error)
	FindByComicIdAndNumber(comicId string, number int) (*entity.Chapter, error)
	Count(id int) (int64, error)
	Delete(id int) error
	CountTotalChapter() (int64, error)
}
