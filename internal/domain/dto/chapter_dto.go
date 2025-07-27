package dto

import "time"

type ChapterAddRequest struct {
	ComicId string `validate:"required,max=36" json:"comic_id"`
	Number  int    `validate:"required" json:"number"`
}
type ChapterUpdateRequest struct {
	Number  *int  `validate:"omitempty" json:"number"`
	Publish *bool `validate:"omitempty" json:"publish"`
}
type ChapterResponse struct {
	Id        int64          `json:"id"`
	ComicId   string         `json:"comic_id"`
	Number    int            `json:"number"`
	Publish   bool           `json:"publish"`
	CreatedAt time.Time      `json:"created_at"`
	Comic     *ComicResponse `json:"comic"`
}
