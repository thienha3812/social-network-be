package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var PRIVATE_KEY = []byte("secret")

func DecodeToken(token string,key string )(interface{},error){
	t,err :=jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
		return PRIVATE_KEY,nil
	})
	if err != nil {
		return PRIVATE_KEY,nil
	}
	if t.Valid {
		value :=  t.Claims.(jwt.MapClaims)[key]
		if value == nil {
			return nil,errors.New("Key không tồn tại")
		}
		return t.Claims.(jwt.MapClaims)[key],nil
	}
	return nil,nil
}