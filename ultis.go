package main

import (
	"errors"
	fmt "fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var PRIVATE_KEY = []byte("secret")

func DecodeToken(c echo.Context, key string) (interface{}, error) {
	token, err := c.Cookie("token")
	t, err := jwt.Parse(token.Value, func(token *jwt.Token) (i interface{}, err error) {
		return PRIVATE_KEY, nil
	})
	if err != nil {
		return PRIVATE_KEY, nil
	}
	if !t.Valid {
		fmt.Println("lỗi nè")
		return nil, errors.New("Token expired")
	}
	if t.Valid {
		value := t.Claims.(jwt.MapClaims)[key]
		if value == nil {
			return nil, errors.New("Key không tồn tại")
		}
		fmt.Println(value)
		return t.Claims.(jwt.MapClaims)[key], nil
	}
	return nil, nil
}
