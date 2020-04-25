package main

import (
	fmt "fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct{}

func (*UserController) Posting(c echo.Context) error {

	request := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	account_id, err := DecodeToken(c, "id")
	if err != nil {
		return nil
	}
	err = db.Create(&Post{AccountID: uint(int(account_id.(float64))), Content: request["content"].(string)}).Error
	if err != nil {
		return nil
	}
	response["success"] = true
	response["message"] = "Đăng bài viết thành công"
	return c.JSON(http.StatusOK, response)
}

func (*UserController) AddFriend(c echo.Context) error {
	request := make(map[string]interface{})

	//response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	account_id, _ := DecodeToken(c, "id")
	var friends Friends
	rows := db.Where("user_1 = ? ", account_id).Where("user_2 =  ? ", int(request["id"].(float64))).Find(&friends).RowsAffected
	if rows == 1 {
		db.Unscoped().Where("user_1 = ? ", account_id).Where("user_2 =  ? ", int(request["id"].(float64))).Delete(Friends{})
	} else {
		db.Create(&Friends{User_1: uint(int(account_id.(float64))), User_2: uint(int(request["id"].(float64)))})
	}
	return nil
}

func (*UserController) LoadProfile(c echo.Context) error {

	request := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	var profile Profile
	var friends Friends
	err := db.Where("account_id = ?", request["id"]).Find(&profile).Error
	if err != nil {
		return nil
	}
	account_id, err := DecodeToken(c, "id")
	db.Where("user_1 = ? ", int(request["id"].(float64))).Where("user_2 =  ? ", account_id).Find(&friends)
	db.Where("user_1 = ? ", account_id).Where("user_2 =  ? ", int(request["id"].(float64))).Find(&friends)
	response["success"] = true
	response["profile"] = echo.Map{
		"full_name": profile.FullName,
		"live_in":   profile.LiveIn,
		"avatar":    profile.Avatar,
	}
	if friends.Status == 0 {
		response["request_status"] = "Thêm bạn bè"
		response["request_status_code"] = 0
	}
	if friends.Status == 1 {
		response["request_status"] = "Đang chờ kết bạn"
		response["request_status_code"] = 1
	}
	if friends.Status == 2 {
		response["request_status"] = "Bạn bè"
		response["request_status_code"] = 2
	}
	if int(account_id.(float64)) == int(friends.User_1) {
		response["self_request"] = true
	}
	return c.JSON(http.StatusOK, response)
}

func (*UserController) LoadRequest(c echo.Context) error {
	fmt.Println("load_request")
	response := make(map[string]interface{})

	account_id, _ := DecodeToken(c, "id")
	var friends []Friends
	type Result struct {
		AccountID uint   `json:"id"`
		Full_Name string `json:"full_name"`
		Avatar    string `json:"avatar"`
	}
	var result []Result
	db.Where("status = 1").Where("user_2 = ?", account_id).Find(&friends)
	for _, friend := range friends {
		var currentProfile Profile
		db.Where("account_id = ?", friend.User_1).Select([]string{"account_id", "full_name", "avatar"}).Find(&currentProfile)
		result = append(result, Result{
			AccountID: currentProfile.AccountID,
			Full_Name: currentProfile.FullName,
			Avatar:    currentProfile.Avatar,
		})
	}
	response["list_request"] = result
	response["sucess"] = true
	return c.JSON(http.StatusOK, response)
}

func (*UserController) AcceptFriend(c echo.Context) error {

	request := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	account_id, _ := DecodeToken(c, "id")
	var friends Friends
	db.Model(&friends).Where("user_1 = ? ", int(request["id"].(float64))).Where("user_2 =  ? ", account_id).Update(Friends{Status: 2})
	response["success"] = true
	return c.JSON(http.StatusOK, response)
}

func (*UserController) GetUserOnline(c echo.Context) error {
	fmt.Println("get_useronline")
	request := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	account_id, _ := DecodeToken(c, "id")
	type ResultForFriends struct {
		Current_Friend uint
	}
	var resultForFriends []ResultForFriends
	db.Raw(`SELECT CASE WHEN user_1 <> ? THEN user_1 ELSE user_2 END as current_friend FROM "Friends" WHERE "Friends".user_1 = ? OR "Friends".user_2 = ?`, account_id, account_id, account_id).Scan(&resultForFriends)
	arrIDForListFriend := []uint{}
	// Ngay đây có lỗi khi chưa có bạn bè
	if len(resultForFriends) == 0 {
		response["list_user"] = []int{}
		c.JSON(200, response)
	}
	for i := range resultForFriends {
		arrIDForListFriend = append(arrIDForListFriend, resultForFriends[i].Current_Friend)
	}
	type ResultForProfile struct {
		FullName string `json:"full_name"`
		Avatar   string `json:"avatar"`
		Username string `json:"username"`
		IsOnline uint   `json:"is_online"`
	}
	var resultForProfile []ResultForProfile
	db.Raw(`SELECT "Account".username,"Profile".*,
	CASE
		WHEN "Profile".account_id =  ANY(ARRAY(SELECT "Account_Online".account_id FROM "Account_Online"  WHERE "Account_Online".account_id != ?))
		THEN 1
		ELSE 0
	END AS is_online
	FROM "Profile"
	LEFT JOIN "Account" ON "Account".id = "Profile".account_id
	WHERE "Profile".account_id IN(?)`, account_id, arrIDForListFriend).Scan(&resultForProfile)
	response["list_user"] = resultForProfile
	return c.JSON(200, response)
}
