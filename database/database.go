package database

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Conn() *gorm.DB {
	dsn := "host=192.168.1.28 user=miauw password=password dbname=miauw port=5432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(" [!] Failed to connect to database.")
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Panic(" [!] Failed to open pool.")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

func RedisConn() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.28:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}
