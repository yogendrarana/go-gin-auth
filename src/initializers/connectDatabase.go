package initializers

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDatabase() {
	var err error

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		// Handle the error as needed, e.g., log it or exit the program
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	fmt.Println("Database connected")
}

// GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
