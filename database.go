package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	DB *gorm.DB
}

func (*Database) init() *gorm.DB {
	// db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=test sslmode=disable password=postgres")
	db, err := gorm.Open("postgres", "host=ec2-34-197-141-7.compute-1.amazonaws.com port=5432 user=eeogohfzqibxdd dbname=der1geoiabi8nk password=20da858827b726ecbff23779708cd57f281805ed11722dfd03aaa2d779004cca")
	if err != nil {
		panic(err)
	}
	return db
}
