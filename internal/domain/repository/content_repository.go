package repository

import "webtoon/internal/domain/entity"

type ContentRepository interface {
	Save(content *entity.Content) error
}
