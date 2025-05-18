package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// TODO: move it to config file to be more formal and all username and password should be encrypted
	dsn := "host=postgres port=5432 user=postgres password=yourpassword dbname=wallet_db sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database initialized successfully.")
}
