package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	fmt "fmt"
	"net/http"
)

type AccountController struct{}

func (*AccountController) Signin(c echo.Context) error {
	request := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	var account Account
	rows := db.Where("username = ?", request["username"]).Where("password = ?", request["password"]).Find(&account).RowsAffected
	if rows == 0 {
		response["success"] = false
		response["message"] = "Tài khoản hoặc mật khẩu không đúng"
		return c.JSON(http.StatusOK, response)
	}
	// Query profile table
	var profile Profile
	db.Where("account_id = ?", account.ID).Find(&profile)
	// Set token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = account.Username
	claims["id"] = account.ID
	//
	ss, err := token.SignedString([]byte("secret"))
	fmt.Println(ss)
	if err != nil {
		panic(err)
		return nil
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    ss,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})
	// set response
	response["success"] = true
	response["message"] = "Đăng nhập thành công"
	response["user_infor"] = echo.Map{
		"avatar":    profile.Avatar,
		"full_name": profile.Full_Name,
		"username":  account.Username,
		"id":        account.ID,
	}
	return c.JSON(http.StatusOK, response)
}
