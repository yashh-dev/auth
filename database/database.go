package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func Conn() *gorm.DB {
	dsn := "host=192.168.1.28 user=miauw_user password=miauw_password dbname=miauw port=5432 sslmode=disable TimeZone=Europe/Berlin"
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
