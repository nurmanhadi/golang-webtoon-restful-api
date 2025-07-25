package mysql

import (
	"webtoon/internal/domain/comic"

	"gorm.io/gorm"
)

type comicStorage struct {
	db *gorm.DB
}

func NewComicStorage(db *gorm.DB) comic.ComicRepository {
	return &comicStorage{db: db}
}
func (r *comicStorage) Save(comic *comic.Comic) error {
	return r.db.Save(&comic).Error
}
func (r *comicStorage) FindById(id string) (*comic.Comic, error) {
	var comic *comic.Comic
	err := r.db.Where("id = ?", id).First(&comic).Error
	if err != nil {
		return nil, err
	}
	return comic, nil
}
func (r *comicStorage) FindAll(page int, size int) ([]comic.Comic, error) {
	var comics []comic.Comic
	err := r.db.Offset((page - 1) * size).Limit(size).Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}
func (r *comicStorage) CountTotal() (int64, error) {
	var count int64
	err := r.db.Model(&comic.Comic{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
