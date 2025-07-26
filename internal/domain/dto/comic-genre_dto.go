package dto

type ComicGenreAddRequest struct {
	ComicId string `validate:"required,max=36" json:"comic_id"`
	GenreId int    `validate:"required" json:"genre_id"`
}
