package repository

import "webtoon/internal/domain/entity"

type AuthRepository interface {
	Save(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	CountByUsername(username string) (int64, error)
}
