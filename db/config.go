package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitializeDB initializes the database connection and runs migrations
func InitializeDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=aimbot dbname=mindaroom port=5432 sslmode=disable"
	var err error

	// Open a connection to the database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Run migrations to ensure the tables exist
	// err = DB.AutoMigrate(&models.User{}, &models.Message{}, &models.GroupMember{}, &models.Session{})
	// if err != nil {
	// 	log.Fatalf("Could not migrate database: %v", err)
	// }

	fmt.Println("Database connected successfully")

	return DB
}
