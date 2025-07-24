package comic

import comictype "webtoon/pkg/comic-type"

type ComicAddRequest struct {
	Title    string         `validate:"required,max=100"`
	Synopsis string         `validate:"required"`
	Author   string         `validate:"required,max=50"`
	Artist   string         `validate:"required,max=50"`
	Type     comictype.TYPE `validate:"required,oneof=manga manhua manhwa"`
}
