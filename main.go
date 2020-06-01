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
	db.AutoMigrate(Profile{})
	db.AutoMigrate(Post{})
	e.POST("/api/account/signin", accountController.Signin)
	e.POST("/api/user/posting", userController.Posting)
	e.POST("/api/user/load-profile", userController.LoadProfile)
	e.POST("/api/user/add-friend", userController.AddFriend)
	e.GET("/api/user/load-request", userController.LoadRequest)
	e.POST("/api/user/accept-friend", userController.AcceptFriend)
	e.POST("/api/user/cancle-request", userController.CancleAcceptFriend)
	e.GET("/api/user/useronline", userController.GetUserOnline)
	e.POST("/api/user/signout", userController.Signout)
	//
	e.POST("/api/places/list-places", placesController.ListPlaces)
	e.POST("/api/places/place-by-id", placesController.GetPlaceByID)
	e.POST("/api/places/user-review", placesController.UserReview)
	e.POST("/api/places/add-place", placesController.AddPlace)
	//
	e.POST("/api/post/user-like", postController.LikePost)
	e.POST("/api/post/user-comment", postController.UserComment)

	e.Use(CheckToken)
	e.Logger.Fatal(e.Start(":8080"))
}
