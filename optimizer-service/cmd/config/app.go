// Package config
package config

import (
	"fmt"
	"log"
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
		db, err := connectToPosgress()
		if err != nil {
			fmt.Println("Error connecting to database, retrying")
		} else {
			log.Printf("Connected to database")
			// err = db.AutoMigrate(&models.User{}, &models.PersonalToken{})
			if err != nil {
				fmt.Println("Error migrating the schema")
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

func connectToPosgress() (*gorm.DB, error) {
	DSN := os.Getenv("DSN")
	fmt.Printf("DSN %v\n%", DSN)
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database %v", err)
		return nil, err
	}
	return db, nil
}
