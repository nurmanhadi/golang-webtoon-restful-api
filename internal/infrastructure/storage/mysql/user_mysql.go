package mysql

import (
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"

	"gorm.io/gorm"
)

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) repository.UserRepository {
	return &userStorage{db: db}
}
func (r *userStorage) Save(user *entity.User) error {
	return r.db.Save(&user).Error
}
func (r *userStorage) FindById(id string) (*entity.User, error) {
	var user *entity.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *userStorage) CountTotalUser() (int64, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("role = ?", "user").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
