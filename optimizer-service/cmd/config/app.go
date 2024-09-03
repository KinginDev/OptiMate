// Package config
package config

import (
	"log"
	"optimizer-service/cmd/internal/models"
	"optimizer-service/cmd/internal/storage"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB      *gorm.DB
	Storage storage.Storage
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
func (app *Config) InitStorage() storage.Storage {
	app.Storage = setUpStorage()
	return app.Storage
}

func connectToPostgress() (*gorm.DB, error) {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	log.Printf("DATABASE_URL %v\n", DATABASE_URL)
	db, err := gorm.Open(postgres.Open(DATABASE_URL), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database %v", err)
		return nil, err
	}
	return db, nil
}

// setUpStorage godoc
// @Summary Set up storage
// @Description Set up storage
func setUpStorage() storage.Storage {
	//Get disk from env
	disk := os.Getenv("DISK")
	switch disk {
	case "local":
		//setup local
		basePath := "./storage/uploads"

		// checks if the directory exists.
		_, err := os.Stat(basePath)

		// If it doesn't, os.IsNotExist(err) returns`true`
		if os.IsNotExist(err) {
			// Create the directory
			errDir := os.MkdirAll(basePath, 0755)
			if errDir != nil {
				log.Printf("Error creating directory %v", errDir)
			}
		} else if err != nil {
			// If os.Stat returned an error other than ErrNotExist, handle it
			log.Printf("Error checking directory %v", err)
		}

		return storage.NewLocalStorage(basePath)
	case "minio":
		//setup minio
		m := &MinioConfig{}
		endpoint := os.Getenv("MINIO_ENDPOINT")
		rootUser := os.Getenv("MINIO_ROOT_USER")
		rootPassword := os.Getenv("MINIO_ROOT_PASSWORD")
		useSSL := m.GetUseSSL() // Configurable based on your setup

		c := NewMinioClient(
			endpoint,
			rootUser,
			rootPassword,
			useSSL,
		)
		log.Printf("Minio Config is %v\n", c)
		bucketName := "optimate"
		return storage.NewMinIOStorage(c, bucketName)
	default:
		log.Println("Unsupported Storage Type")
	}

	return nil
}
