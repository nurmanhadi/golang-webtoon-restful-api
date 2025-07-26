package dto

import "webtoon/pkg"

type GenreAddRequest struct {
	Name string `validate:"required,max=50" json:"name"`
}
type GenreResponse struct {
	Id     int                          `json:"id"`
	Name   string                       `json:"name"`
	Comics *pkg.Paging[[]ComicResponse] `json:"comics"`
}
