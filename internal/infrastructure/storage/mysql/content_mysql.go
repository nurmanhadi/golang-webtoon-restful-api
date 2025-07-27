package mysql

import (
	"webtoon/internal/domain/entity"
	"webtoon/internal/domain/repository"

	"gorm.io/gorm"
)

type contentStorage struct {
	db *gorm.DB
}

func NewContentStorage(db *gorm.DB) repository.ContentRepository {
	return &contentStorage{db: db}
}
func (r *contentStorage) Save(content *entity.Content) error {
	return r.db.Save(content).Error
}
