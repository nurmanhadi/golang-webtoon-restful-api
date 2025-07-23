package user

type User struct {
	Id             string `gorm:"type:varchar(36);primaryKey"`
	Username       string `gorm:"type:varchar(100);unique;not null"`
	Password       string `gorm:"type:varchar(100);not null"`
	Role           string `gorm:"type:enum('admin','user');not null;default 'user'"`
	AvatarFilename string `gorm:"type:varchar(100)"`
	AvatarUrl      string `gorm:"type:varchar(255)"`
}
