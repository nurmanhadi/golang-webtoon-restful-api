package auth

import (
	"webtoon/internal/domain/user"
)

type AuthRepository interface {
	Save(user *user.User) error
	FindByUsername(username string) (*user.User, error)
	CountByUsername(username string) (int64, error)
}
