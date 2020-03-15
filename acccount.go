package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"regexp"
	"github.com/dgrijalva/jwt-go"
)

type AccountController struct{}

func (*AccountController) loginAccount(c *gin.Context) {
	m := make(map[string]interface{})
	response := make(map[string]interface{})
	var account Account
	if err := c.Bind(&m); err != nil {
		c.Next()
		return
	}
	fmt.Println(m)
	//
	err := db.Model(&account).Where("username = ?",m["username"]).Where("password = ?",m["password"]).Select()
	if err != nil {
		response["code"] = 101
		response["message"] = "Đăng nhập thất bại"
		c.JSON(200, response)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256,jwt.MapClaims{
		"username" : account.Username,
	})
	tokenString, _ := token.SigningString()
	fmt.Println(tokenString)
	response["code"] = 102
	response["message"] = "Đăng nhập thành công"
	_data := make(map[string]interface{})
	_data["token"] = tokenString
	response["data"] = _data
	c.JSON(200, response)
}
func (*AccountController) forgotPassword(c *gin.Context) {
	m := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&m); err != nil {
		response["code"] = 101 // Dữ liệu không hợp lệ
		response["message"] = "Dữ liệu không hợp lệ"
		c.JSON(200, response)
		return
	}

}

func (*AccountController) register(c *gin.Context) {
	m := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&m); err != nil {
		response["code"] = 101 // Dữ liệu không hợp lệ
		response["message"] = "Dữ liệu không hợp lệ"
		c.JSON(200, response)
		return
	}
	format := []string{"username", "password", "email"}
	for i := range format {
		if m[format[i]].(string) == "" {
			response["code"] = 101 // Dữ liệu không hợp lệ
			response["message"] = "Dữ liệu không hợp lệ"
			c.JSON(200, response)
			return
		}
	}
	// Validate
	reEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	reUsername := regexp.MustCompile("^[a-z0-9]{6,20}$")
	if ok := reUsername.MatchString(m["username"].(string)); !ok {
		response["code"] = 101 // Dữ liệu không hợp lệ
		response["message"] = "Tài khoản  không hợp lệ"
		c.JSON(200, response)
		return
	}
	if ok := reEmail.MatchString(m["email"].(string)); !ok {
		response["code"] = 101 // Dữ liệu không hợp lệ
		response["message"] = "Email không hợp lệ"
		c.JSON(200, response)
		return
	}

	account := Account{
		Username: m["username"].(string),
		Password: m["password"].(string),
		Email:    m["email"].(string),
	}
	//Thêm account vào table
	if err := db.Insert(&account); err != nil {
		response["code"] = 102 // Tài khoản hoặc email đã được đăng ký
		response["message"] = " Tài khoản hoặc email đã được đăng ký"
		c.JSON(200, response)
		return
	}
	// Send email kích hoạt tài khoản
	//auth := smtp.PlainAuth("", "thienhamaimai1@gmail.com", "fupkijjadpvatryi", "smtp.gmail.com")
	_uuid, _ := uuid.NewUUID()
	if err := redisClient.Set(_uuid.String(),m["username"].(string),0).Err(); err != nil {
		c.Next()
		return
	}
	// Send mail
	to := []string{m["email"].(string)}
	content := fmt.Sprintf("Cảm ơn bạn đã đăng ký và sử dụng dịch vụ của chúng tôi bạn vui lòng nhấp vào link dưới đây để kích hoạt tài khoản  " +
		"\r\n" +
		"https://localhost:3000/active/%s",_uuid.String())
	msg := []byte(
		"Subject: Kích hoạt tài khoản" +
			"\r \n" +
			content,
	)
	mail := Mail{
		Msg:msg,
		To:to,
	}
	go sendMail(chanEmail)
	chanEmail <- mail
	response["code"] = 103 // Đăng ký tài khoản thành công
	response["message"] = "Đăng ký tài khoản thành công"
	c.JSON(200, response)

}
func (*AccountController) changePassword(c *gin.Context) {
	m := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&m); err != nil {
		c.Next()
	}
	var account Account
	err := db.Model(&account).Where("username = (?)", m["username"].(string)).Select()
	if err != nil {
		c.Next()
	}
	account.Password = "123"
	if err := db.Update(&account); err != nil {
		c.Next()
	}
	response["code"] = 102
	response["message"] = "Đăng nhập thành công"
	c.JSON(200, response)
}
