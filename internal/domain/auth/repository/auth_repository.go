package repository

import (
	"webtoon/internal/domain/user/entity"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Save(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	CountByUsername(username string) (int64, error)
}
type repository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &repository{db: db}
}

func (r *repository) Save(user *entity.User) error {
	return r.db.Save(&user).Error
}
func (r *repository) FindByUsername(username string) (*entity.User, error) {
	var user *entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *repository) CountByUsername(username string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
