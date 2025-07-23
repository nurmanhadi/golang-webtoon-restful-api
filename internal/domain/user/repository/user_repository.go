package repository

import "gorm.io/gorm"

type UserRepository interface{}

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) {}
