package entity

type ComicGenre struct {
	Id      int64  `gorm:"primaryKey;type:bigint"`
	ComicId string `gorm:"type:varchar(36);not null"`
	GenreId int    `gorm:"type:int;not null"`
	Comic   *Comic
	Genre   *Genre
}
