package config

import (
	"log"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	validator := validator.New()
	log.Println("validator activate")
	return validator
}
