package entity

type Genre struct {
	Id         int    `gorm:"primaryKey;type:bigint"`
	Name       string `gorm:"type:varchar(50);not null"`
	ComicGenre []ComicGenre
}
