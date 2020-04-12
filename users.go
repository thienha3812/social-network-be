
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserController struct {}




func (*UserController) Posting(c echo.Context) error {
	token, _ := c.Cookie("token")
	request := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	account_id,err  := DecodeToken(token.Value,"id")
	if err != nil {
		return nil
	}
	err = db.Create(&Post{AccountID:uint(int(account_id.(float64))),Content:request["content"].(string)}).Error
	if err != nil {
		return nil
	}
	response["success"] = true
	response["message"] = "Đăng bài viết thành công"
	return c.JSON(http.StatusOK, response)
}

func (*UserController) AddFriend(c echo.Context) error {
	request := make(map[string]interface{})
	token, _ := c.Cookie("token")
	//response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	account_id , _ := DecodeToken(token.Value,"id")
	var friends Friends
	rows := db.Where("user_1 = ? ",account_id).Where("user_2 =  ? ",int(request["id"].(float64))).Find(&friends).RowsAffected
	if rows == 1 {
		db.Unscoped().Where("user_1 = ? ",account_id).Where("user_2 =  ? ",int(request["id"].(float64))).Delete(Friends{})
	}else{
		db.Create(&Friends{User_1:uint(int(account_id.(float64))),User_2:uint(int(request["id"].(float64)))})
	}
	return nil
}

func (*UserController) LoadProfile(c echo.Context) error {
	token, _ := c.Cookie("token")
	request := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	var profile Profile
	var friends Friends
	err := db.Where("account_id = ?",request["id"]).Find(&profile).Error
	if err != nil {
		return nil
	}
	account_id,err  := DecodeToken(token.Value,"id")
	db.Where("user_1 = ? ",account_id).Where("user_2 =  ? ",int(request["id"].(float64))).Find(&friends)
	response["success"] = true
	response["profile"] = echo.Map{
		"full_name" : profile.Full_Name,
		"live_in" : profile.Live_In,
		"avatar" : profile.Avatar,

	}
	if friends.Status == 0 {
		response["request_status"] = "Thêm bạn bè"
	}
	if friends.Status == 1 {
		response["request_status"] = "Huỷ lời mời"
	}
	if friends.Status == 2 {
		response["request_status"] = "Bạn bè"
	}
	return c.JSON(http.StatusOK,response)
}

func (*UserController) LoadRequest(c echo.Context) error{
	response := make(map[string]interface{})
	token, _ := c.Cookie("token")
	account_id,_  := DecodeToken(token.Value,"id")
	var friends []Friends
	type Result struct {
		AccountID  uint `json:"id"`
		Full_Name string `json:"full_name"`
	}
	var result []Result
	db.Where("status = 1").Where("user_2 = ?",account_id).Find(&friends)
	for _,friend := range friends  {
		var currentProfile Profile
		db.Where("account_id = ?",friend.User_1).Select([]string{"account_id","full_name"}).Find(&currentProfile)
		result = append(result,Result{
			AccountID: currentProfile.AccountID,
			Full_Name: currentProfile.Full_Name,
		})
	}
	response["list_request"] = result
	response["sucess"] = true
	return c.JSON(http.StatusOK,response)
}
