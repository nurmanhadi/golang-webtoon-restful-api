package repository

import "webtoon/internal/domain/entity"

type ComicGenreRepository interface {
	Save(comicGenre *entity.ComicGenre) error
	DeleteByComicIdAndGenreId(comicId string, genreId int) error
}
