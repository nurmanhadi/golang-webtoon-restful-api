package mysql

import (
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"

	"gorm.io/gorm"
)

type chapterStorage struct {
	db *gorm.DB
}

func NewChapterStorage(db *gorm.DB) repository.ChapterRepository {
	return &chapterStorage{db: db}
}
func (r *chapterStorage) Save(chapter *entity.Chapter) error {
	return r.db.Save(chapter).Error
}
func (r *chapterStorage) FindById(id int) (*entity.Chapter, error) {
	var chapter *entity.Chapter
	err := r.db.Where("id = ?", id).Preload("Contents").First(&chapter).Error
	if err != nil {
		return nil, err
	}
	return chapter, nil
}
func (r *chapterStorage) FindByComicIdAndNumber(comicId string, number int) (*entity.Chapter, error) {
	var chapter *entity.Chapter
	err := r.db.Where("comic_id = ? AND number = ?", comicId, number).Preload("Comic.Chapters").Preload("Contents").First(&chapter).Error
	if err != nil {
		return nil, err
	}
	return chapter, nil
}
func (r *chapterStorage) Count(id int) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Chapter{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *chapterStorage) Delete(id int) error {
	return r.db.Where("id = ?", id).Delete(&entity.Chapter{}).Error
}
