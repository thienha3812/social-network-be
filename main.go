package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/go-pg/pg/v9"
	"github.com/go-redis/redis/v7"
)


func initDatabase() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "admin123",
		Addr:     CONNECT_STRING,
		Database: "network-social",
		OnConnect: func(conn *pg.Conn) error {
			fmt.Println("Connected to database")
			return nil
		},
	})
	return db
}

var db = initDatabase()
var CONNECT_STRING = "database-1.crbkggicegau.us-east-1.rds.amazonaws.com:5432"
var accountController AccountController
var redisClient = redis.NewClient(&redis.Options{DB: 0, Addr: "localhost:6379", Password: ""})
/// Chanel

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowHeaders: []string{"Content-Type"},
	}))
	// Acount Controller

	fmt.Println("----------------")
	//cmd := exec.Command("python3","test.py","AN Cafe, Đường Lê Đại Hành, Thành phố Nha Trang, Khánh Hòa")
	//output , _ := cmd.Output()
	//fmt.Println(string(output))
	r.POST("/account/login", accountController.loginAccount)
	r.POST("/account/changepassword", accountController.changePassword)
	r.POST("/account/register", accountController.register)
	r.Run()
}
