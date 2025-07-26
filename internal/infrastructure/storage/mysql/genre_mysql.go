package mysql

import (
	"webtoon/internal/domain/genre"

	"gorm.io/gorm"
)

type genreStorage struct {
	db *gorm.DB
}

func NewGenreStorage(db *gorm.DB) genre.GenreRepository {
	return &genreStorage{db: db}
}
func (r *genreStorage) Save(genre *genre.Genre) error {
	return r.db.Save(genre).Error
}
func (r *genreStorage) Remove(id int) error {
	return r.db.Where("id = ?", id).Delete(&genre.Genre{}).Error
}
func (r *genreStorage) FindAll() ([]genre.Genre, error) {
	var genres []genre.Genre
	err := r.db.Find(&genres).Error
	if err != nil {
		return nil, err
	}
	return genres, nil
}
func (r *genreStorage) Count(id int) (int64, error) {
	var count int64
	err := r.db.Model(&genre.Genre{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
