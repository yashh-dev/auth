package database

import (
	"log"
	"miauw.social/auth/config"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Conn() *gorm.DB {

	dsn := config.GetConfig().DBUrl
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
	cfg := config.GetConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHOST,
		Password: cfg.RedisPass,
		DB:       0,
	})
	return rdb
}
