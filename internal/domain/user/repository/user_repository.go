package repository

import (
	"webtoon/internal/domain/user/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user *entity.User) error
	FindById(id string) (*entity.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &repository{db: db}
}
func (r *repository) Save(user *entity.User) error {
	return r.db.Save(&user).Error
}
func (r *repository) FindById(id string) (*entity.User, error) {
	var user *entity.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
