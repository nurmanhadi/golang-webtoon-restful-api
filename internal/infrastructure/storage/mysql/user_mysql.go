package mysql

import (
	"webtoon/internal/domain/user"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) Save(user *user.User) error {
	return r.db.Save(&user).Error
}
func (r *userRepository) FindById(id string) (*user.User, error) {
	var user *user.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
