package mysql

import (
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"

	"gorm.io/gorm"
)

type authStorage struct {
	db *gorm.DB
}

func NewAuthStorage(db *gorm.DB) repository.AuthRepository {
	return &authStorage{db: db}
}

func (r *authStorage) Save(user *entity.User) error {
	return r.db.Save(&user).Error
}
func (r *authStorage) FindByUsername(username string) (*entity.User, error) {
	var user *entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *authStorage) CountByUsername(username string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
