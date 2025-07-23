package mysql

import (
	"webtoon/internal/domain/auth"
	"webtoon/internal/domain/user"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) auth.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Save(user *user.User) error {
	return r.db.Save(&user).Error
}
func (r *authRepository) FindByUsername(username string) (*user.User, error) {
	var user *user.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *authRepository) CountByUsername(username string) (int64, error) {
	var count int64
	err := r.db.Model(&user.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
