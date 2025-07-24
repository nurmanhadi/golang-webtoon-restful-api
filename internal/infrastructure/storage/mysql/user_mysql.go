package mysql

import (
	"webtoon/internal/domain/user"

	"gorm.io/gorm"
)

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) user.UserRepository {
	return &userStorage{db: db}
}
func (r *userStorage) Save(user *user.User) error {
	return r.db.Save(&user).Error
}
func (r *userStorage) FindById(id string) (*user.User, error) {
	var user *user.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
