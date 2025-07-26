package repository

import "webtoon/internal/domain/entity"

type ComicGenreRepository interface {
	Save(comicGenre *entity.ComicGenre) error
	DeleteByComicIdAndGenreId(comicId string, genreId int) error
	FindAllByGenreId(genreId int, page int, size int) ([]entity.ComicGenre, error)
	CountByGenreId(genreId int) (int64, error)
}
