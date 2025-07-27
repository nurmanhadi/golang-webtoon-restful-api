package entity

import "time"

type Chapter struct {
	Id        int64  `gorm:"primaryKey;type:bigint"`
	ComicId   string `gorm:"type:varchar(36);not null"`
	Number    int    `gorm:"type:int;not null"`
	Publish   bool   `gorm:"type:bool;not null"`
	CreatedAt time.Time
	Comic     *Comic
}
