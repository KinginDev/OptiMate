package config

import (
	"fmt"
	"log"
	"os"
	"time"
	"user-service/data"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Database *gorm.DB
	Models   data.Models
}

var counts int64

func (c *Config) InitDB() *gorm.DB {

	for {
		db, err := connectToPostgres()
		if err != nil {
			fmt.Println("Error connecting to database, retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			counts++
		} else {
			log.Printf("Connected to database")
			db.AutoMigrate(&data.User{}, &data.PersonalToken{})
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

func connectToPostgres() (*gorm.DB, error) {
	DSN := os.Getenv("DSN")
	fmt.Printf("DSN: %v\n", DSN)
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Printf("failed to connect database %v", err)
		return nil, err
	}
	return db, nil
}
