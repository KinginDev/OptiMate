// Package config
package config

import (
	"log"
	"optimizer-service/cmd/internal/models"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB *gorm.DB
}

func NewConfig() *Config {
	return &Config{}
}

var counts int64

func (app *Config) InitDB() *gorm.DB {
	for {
		db, err := connectToPostgress()
		if err != nil {
			log.Println("Error connecting to database, retrying")
			time.Sleep(5 * time.Second)
			counts++
		} else {
			log.Printf("Connected to database")
			err = db.AutoMigrate(&models.File{}, &models.OptimizationSettings{})
			if err != nil {
				log.Println("Error migrating the schema")
				return nil
			}
			return db
		}

		if counts > 10 {
			log.Printf("Failed to connect to database after 10 retries")
			return nil
		}

		log.Printf("Retrying to connect to database")
		time.Sleep(2 * time.Second)
		continue
	}
}

func connectToPostgress() (*gorm.DB, error) {
	DSN := os.Getenv("DSN")
	log.Printf("DSN %v\n", DSN)
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database %v", err)
		return nil, err
	}
	return db, nil
}
