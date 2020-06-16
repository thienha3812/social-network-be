package main

import (
	"encoding/json"
	"errors"
	fmt "fmt"
	"math/rand"

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
func ConvertToJson(value interface{}) []byte {
	s := fmt.Sprintf("%+v", value)
	b, _ := json.Marshal(s)
	jsonString := string(b)
	fmt.Println(jsonString[1 : len(jsonString)-1])
	return []byte(jsonString[1 : len(jsonString)-1])
}

func PickRandom(list []CustomPlaces) (CustomPlaces, []CustomPlaces) {
	random := rand.Intn(len(list))
	restList := append(list[:random], list[random+1:]...)
	return list[random], restList
}
