package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	DB *gorm.DB
}



func (*Database) init() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=test sslmode=disable password=postgres")
	if err != nil {
		panic(err)
	}
	return db
}