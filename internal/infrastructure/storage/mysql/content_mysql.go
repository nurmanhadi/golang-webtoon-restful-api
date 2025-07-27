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
func (r *contentStorage) FindById(id int) (*entity.Content, error) {
	var content *entity.Content
	err := r.db.Where("id = ?", id).First(&content).Error
	if err != nil {
		return nil, err
	}
	return content, nil
}
func (r *contentStorage) Delete(id int) error {
	return r.db.Where("id = ?", id).Delete(&entity.Content{}).Error
}
