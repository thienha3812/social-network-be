package main

import (
	"log"
	"net/smtp"
)
var chanEmail = make(chan Mail,1000)
type Mail struct {
	To []string
	Msg []byte
}
func sendMail( email <-chan Mail){
	auth := smtp.PlainAuth("", "thienhamaimai1@gmail.com", "", "smtp.gmail.com")
	for e := range email {
		if err := smtp.SendMail("smtp.gmail.com:587", auth, "thienhamaimai1@gmail.com", e.To, e.Msg); err != nil {
			log.Fatal(err)
		}
	}
}
