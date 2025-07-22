package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysql() *gorm.DB {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_MYSQL_URL")), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf("mysql error: %s", err.Error())
	}
	log.Println("connected to database")

	pool, err := db.DB()
	if err != nil {
		log.Fatalf("database pooling error: %s", err.Error())
	}
	idleConns, err := strconv.ParseInt(os.Getenv("DB_POOL_MAX_IDLE_CONNS"), 10, 32)
	if err != nil {
		log.Fatalf("parse string to int idleConns error: %s", err.Error())
	}
	openConns, err := strconv.ParseInt(os.Getenv("DB_POOL_MAX_OPEN_CONNS"), 10, 32)
	if err != nil {
		log.Fatalf("parse string to int openConns error: %s", err.Error())
	}
	idleTime, err := strconv.ParseInt(os.Getenv("DB_POOL_MAX_IDLE_TIME"), 10, 32)
	if err != nil {
		log.Fatalf("parse string to int idleTime error: %s", err.Error())
	}
	lifetime, err := strconv.ParseInt(os.Getenv("DB_POOL_MAX_LIFETIME"), 10, 32)
	if err != nil {
		log.Fatalf("parse string to int lifetime error: %s", err.Error())
	}
	pool.SetMaxIdleConns(int(idleConns))
	pool.SetMaxOpenConns(int(openConns))
	pool.SetConnMaxIdleTime(time.Duration(idleTime) * time.Minute)
	pool.SetConnMaxLifetime(time.Duration(lifetime) * time.Minute)

	log.Printf("db pooling max idle conns| %d", idleConns)
	log.Printf("db pooling max open conns| %d", openConns)
	log.Printf("db pooling max idle time | %d minute", idleTime)
	log.Printf("db pooling max lifetime  | %d minute", lifetime)
	return db
}
