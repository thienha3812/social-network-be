package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// init struct here
var database Database
var accountController AccountController
var userController UserController
var placesController PlacesController
var postController PostController
var messageController MessageController

// init all necessary
var db = database.init()
var localServer = "http://localhost:8080"

func main() {
	e := echo.New()
	e.Static("/public", "./assets/")
	db.LogMode(true)
	e.Debug = true
	db.AutoMigrate(Places{})
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Set-Cookie"},
	}))
	//
	db.AutoMigrate(Account{})
	db.AutoMigrate(Friends{})
	db.AutoMigrate(Profile{})
	db.AutoMigrate(Post{})
	db.AutoMigrate(Message{})
	db.AutoMigrate(Conversation{})
	db.AutoMigrate(AccountOnline{})
	db.AutoMigrate(Images{})
	db.AutoMigrate(Comment{})
	db.AutoMigrate(Places{})
	//
	e.POST("/api/account/signin", accountController.Signin)
	e.POST("/api/user/posting", userController.Posting)
	e.POST("/api/user/load-profile", userController.LoadProfile)
	e.POST("/api/user/add-friend", userController.AddFriend)
	e.GET("/api/user/load-request", userController.LoadRequest)
	e.POST("/api/user/accept-friend", userController.AcceptFriend)
	e.POST("/api/user/cancle-request", userController.CancleAcceptFriend)
	e.GET("/api/user/user-online", userController.GetUserOnline)
	e.POST("/api/user/signout", userController.Signout)
	//Places controller
	e.POST("/api/places/list-places", placesController.ListPlaces)
	e.POST("/api/places/place-by-id", placesController.GetPlaceByID)
	e.POST("/api/places/user-review", placesController.UserReview)
	e.POST("/api/places/add-place", placesController.AddPlace)
	e.GET("/api/places/get-place-for-index-page", placesController.GetPlacesForIndexPage)
	//
	e.POST("/api/post/user-like", postController.LikePost)
	e.POST("/api/post/user-comment", postController.UserComment)

	//
	e.POST("/api/message/get-history", messageController.GetHistoryMessage)
	e.Use(CheckToken)
	e.Logger.Fatal(e.Start(":8080"))
}
