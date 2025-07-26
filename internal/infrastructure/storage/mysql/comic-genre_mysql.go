package mysql

import (
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"

	"gorm.io/gorm"
)

type comicGenreStorage struct {
	db *gorm.DB
}

func NewComicGenreStorage(db *gorm.DB) repository.ComicGenreRepository {
	return &comicGenreStorage{db: db}
}
func (r *comicGenreStorage) Save(comicGenre *entity.ComicGenre) error {
	return r.db.Save(comicGenre).Error
}
func (r *comicGenreStorage) FindAllByGenreId(genreId int, page int, size int) ([]entity.ComicGenre, error) {
	var comicGenre []entity.ComicGenre
	err := r.db.
		Offset((page-1)*size).
		Limit(size).
		Where("genre_id", genreId).
		Preload("Comic").
		Find(&comicGenre).Error
	if err != nil {
		return nil, err
	}
	return comicGenre, nil
}
func (r *comicGenreStorage) CountByGenreId(genreId int) (int64, error) {
	var count int64
	err := r.db.Model(&entity.ComicGenre{}).Where("genre_id = ?", genreId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}
func (r *comicGenreStorage) DeleteByComicIdAndGenreId(comicId string, genreId int) error {
	return r.db.Where("comic_id = ? AND genre_id = ?", comicId, genreId).Delete(&entity.ComicGenre{}).Error
}
