package comic

import (
	"time"
	comictype "webtoon/pkg/comic-type"
)

type ComicAddRequest struct {
	Title    string         `validate:"required,max=100"`
	Synopsis string         `validate:"required"`
	Author   string         `validate:"required,max=50"`
	Artist   string         `validate:"required,max=50"`
	Type     comictype.TYPE `validate:"required,oneof=manga manhua manhwa"`
}
type ComicUpdateRequest struct {
	Title    string         `validate:"omitempty,max=100"`
	Synopsis string         `validate:"omitempty"`
	Author   string         `validate:"omitempty,max=50"`
	Artist   string         `validate:"omitempty,max=50"`
	Type     comictype.TYPE `validate:"omitempty,oneof=manga manhua manhwa"`
}
type ComicResponse struct {
	Id            string         `json:"id"`
	Title         string         `json:"title"`
	Synopsis      string         `json:"synopsis"`
	Author        string         `json:"author"`
	Artist        string         `json:"artist"`
	Type          comictype.TYPE `json:"type"`
	CoverFilename string         `json:"cover_filename"`
	CoverUrl      string         `json:"cover_url"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
