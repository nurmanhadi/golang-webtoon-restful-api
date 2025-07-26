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
func (r *comicGenreStorage) DeleteByComicIdAndGenreId(comicId string, genreId int) error {
	return r.db.Where("comic_id = ? AND genre_id = ?", comicId, genreId).Delete(&entity.ComicGenre{}).Error
}
