package db

import (
	"fmt"
	"midterm/env"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DbConnect *gorm.DB

func ConnectDB() (*gorm.DB, error) {

	connStr := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", env.Hostname, env.Name, env.DB, env.Password, env.Port)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	DbConnect = db

	return DbConnect, nil
}
