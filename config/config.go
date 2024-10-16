package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := "experiment:Experiment@123@tcp(127.0.0.1:3306)/experiment_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	log.Println("Database connection established.")
	return db
}
