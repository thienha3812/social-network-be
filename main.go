package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)
// init struct here
var database Database
var accountController AccountController
var userController UserController

// init all necessary
var db = database.init()

func main() {

	e := echo.New()
	e.Debug = true
	db.AutoMigrate(Friends{})
	//e.Use(CheckToken)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowCredentials:true,
		AllowHeaders: []string{ "Content-Type","Set-Cookie"},
	}))
	e.POST("/api/account/signin",accountController.Signin)
	//
	e.POST("/api/user/posting",userController.Posting)
	e.POST("/api/user/loadprofile",userController.LoadProfile)
	e.POST("/api/user/addfriend",userController.AddFriend)
	e.GET("/api/user/loadrequest",userController.LoadRequest)
	e.Logger.Fatal(e.Start(":8080"))
}