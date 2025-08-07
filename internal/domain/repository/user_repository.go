package repository

import "webtoon/internal/domain/entity"

type UserRepository interface {
	Save(user *entity.User) error
	FindById(id string) (*entity.User, error)
	CountTotalUser() (int64, error)
}
