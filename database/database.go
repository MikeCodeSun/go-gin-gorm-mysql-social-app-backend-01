package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)



func ConnectDB() *gorm.DB{
	// .env varible uri, name, password, username
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error laoding .env file")
	}
	uri := os.Getenv("DB_URI")
	db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_username := os.Getenv("DB_USERNAME")
	// format dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_password, uri, db_name)
	// open mysql 
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatal("Error Db Open")
	}

	return db
}

