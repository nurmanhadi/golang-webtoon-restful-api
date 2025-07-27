package entity

type Content struct {
	Id        int64  `gorm:"primaryKey;type:bigint"`
	ChapterId int64  `gorm:"type:bigint;not null"`
	Filename  string `gorm:"type:varchar(100); not null"`
	Url       string `gorm:"type:varchar(255);not null"`
	Chapter   *Chapter
}
