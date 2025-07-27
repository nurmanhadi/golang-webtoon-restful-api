package mysql

import (
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"

	"gorm.io/gorm"
)

type comicStorage struct {
	db *gorm.DB
}

func NewComicStorage(db *gorm.DB) repository.ComicRepository {
	return &comicStorage{db: db}
}
func (r *comicStorage) Save(comic *entity.Comic) error {
	return r.db.Save(&comic).Error
}
func (r *comicStorage) FindById(id string) (*entity.Comic, error) {
	var comic *entity.Comic
	err := r.db.Where("id = ?", id).Preload("ComicGenre.Genre").Preload("Chapters").First(&comic).Error
	if err != nil {
		return nil, err
	}
	return comic, nil
}
func (r *comicStorage) FindAll(page int, size int) ([]entity.Comic, error) {
	var comics []entity.Comic
	err := r.db.Offset((page - 1) * size).Limit(size).Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}
func (r *comicStorage) CountTotal() (int64, error) {
	var count int64
	err := r.db.Model(&entity.Comic{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *comicStorage) CountById(id string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Comic{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *comicStorage) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&entity.Comic{}).Error
}
func (r *comicStorage) Search(key string, page int, size int) ([]entity.Comic, error) {
	var comics []entity.Comic
	keyword := "%" + key + "%"
	err := r.db.
		Offset((page-1)*size).
		Limit(size).
		Where("title  LIKE ? OR author LIKE ? OR artist LIKE ?", keyword, keyword, keyword).
		Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}
func (r *comicStorage) CountTotalByKeyword(key string) (int64, error) {
	var count int64
	keyword := "%" + key + "%"
	err := r.db.
		Model(&entity.Comic{}).
		Where("title  LIKE ? OR author LIKE ? OR artist LIKE ?", keyword, keyword, keyword).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
