package entity

import (
	"time"
	comictype "webtoon/pkg/comic-type"
)

type Comic struct {
	Id            string         `gorm:"type:varchar(36);not null"`
	Title         string         `gorm:"type:varchar(100);not null"`
	Synopsis      string         `gorm:"type:text;not null"`
	Author        string         `gorm:"type:varchar(50);not null"`
	Artist        string         `gorm:"type:varchar(50);not null"`
	Type          comictype.TYPE `gorm:"type:enum('manga','manhua','manhwa');not null"`
	CoverFilename string         `gorm:"type:varchar(100);not null"`
	CoverUrl      string         `gorm:"type:varchar(255);not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ComicGenre    []ComicGenre
}
