package mysql

import (
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"

	"gorm.io/gorm"
)

type genreStorage struct {
	db *gorm.DB
}

func NewGenreStorage(db *gorm.DB) repository.GenreRepository {
	return &genreStorage{db: db}
}
func (r *genreStorage) Save(genre *entity.Genre) error {
	return r.db.Save(genre).Error
}
func (r *genreStorage) Remove(id int) error {
	return r.db.Where("id = ?", id).Delete(&entity.Genre{}).Error
}
func (r *genreStorage) FindAll() ([]entity.Genre, error) {
	var genres []entity.Genre
	err := r.db.Find(&genres).Error
	if err != nil {
		return nil, err
	}
	return genres, nil
}
func (r *genreStorage) Count(id int) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Genre{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
