package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"log"
	"refactored-robot/internal/package/models"
)

func Init() *gorm.DB {
	dbURL := "postgres://postgres:postgres@localhost:5432/postgres"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("Connected")
	}

	db.AutoMigrate(&models.User{})

	return db
}
