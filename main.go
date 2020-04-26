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
	db.LogMode(true)
	e.Debug = true
	db.AutoMigrate(Account{})
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Set-Cookie"},
	}))
	//
	db.AutoMigrate(Profile{})
	e.POST("/api/account/signin", accountController.Signin)
	e.POST("/api/user/posting", userController.Posting)
	e.POST("/api/user/loadprofile", userController.LoadProfile)
	e.POST("/api/user/addfriend", userController.AddFriend)
	e.GET("/api/user/loadrequest", userController.LoadRequest)
	e.POST("/api/user/acceptfriend", userController.AcceptFriend)
	e.GET("/api/user/useronline", userController.GetUserOnline)
	e.POST("/api/user/signout", userController.Signout)
	e.Use(CheckToken)
	e.Logger.Fatal(e.Start(":8080"))
}
