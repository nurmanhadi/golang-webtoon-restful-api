package config

import (
	"log"

	"github.com/joho/godotenv"
)

func NewEnvirontment() error {
	godotenv.Load()
	log.Println("environtment activate")
	return nil
}
