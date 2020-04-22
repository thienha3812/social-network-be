package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() != "/api/account/signin" {
			token, err := c.Cookie("token")
			if err != nil {
				return c.String(500, "Token not exist")
			}
			t, _ := jwt.Parse(token.Value, func(token *jwt.Token) (i interface{}, err error) {
				return []byte("secret"), nil
			})
			if !t.Valid {
				return c.String(500, "Token invalid")
			}
		}
		return next(c)
	}
}
