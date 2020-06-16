package main

import (
	fmt "fmt"
	"net/http"

	"github.com/jinzhu/gorm/dialects/postgres"
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
	if err := c.Bind(&request); err != nil {
		return nil
	}
	account_id, _ := DecodeToken(c, "id")
	if rows := db.Model(&Friends{}).Find("user_1 = ? AND user_2 = ?", account_id, request["user_id"]).RowsAffected; rows == 0 {
		if err := db.Create(&Friends{User_1: uint(account_id.(float64)), User_2: uint(request["user_id"].(float64)), Status: 1}).Error; err != nil {
			return c.NoContent(500)
		}
	}
	return c.NoContent(200)
}

func (*UserController) LoadProfile(c echo.Context) error {
	request := make(map[string]interface{})
	response := echo.Map{}
	account_id, _ := DecodeToken(c, "id")
	if err := c.Bind(&request); err != nil {
		return nil
	}
	if account_id == request["account_id"] || request["account_id"] == nil {
		type Result struct {
			CountPost int      `json:"count_post"`
			Images    []Images `json:"images"`
			Posts     []struct {
				ID        uint           `json:"id"`
				Content   string         `json:"content"`
				Avatar    string         `json:"avatar"`
				FullName  string         `json:"full_name"`
				AccountID string         `json:"account_id"`
				UserLiked postgres.Jsonb `json:"user_liked"`
				Comments  postgres.Jsonb `json:"comments"`
				Images    postgres.Jsonb `json:"images"`
			} `json:"posts"`
		}
		var result Result
		query := ` SELECT "Post".*,"Profile".avatar,"Profile".full_name , "Account".username,
		(
				SELECT COALESCE(json_agg(DISTINCT item),'[]'::json)
				FROM (
						SELECT "Images".* FROM "Images" WHERE "Images".id = ANY("Post".image_ids)
				) item
		) as images,
		(
				SELECT COALESCE(json_agg(DISTINCT item),'[]'::json)
				FROM ( 
						SELECT "Profile".* FROM "Profile" WHERE "Profile".account_id = ANY("Post".account_liked)) item
		) as user_liked ,
		(
				SELECT COALESCE(json_agg(DISTINCT item),'[]'::json)
				FROM (SELECT "Comment".*,"Profile".* as user FROM "Comment","Profile"
				WHERE "Comment".account_id = "Profile".account_id AND "Comment".post_id = "Post".id
				) as item
		) as comments
		FROM "Post","Profile","Places","Account"
		WHERE "Post".account_id = ? AND "Post".account_id = "Profile".account_id AND "Profile".account_id = "Account".id
		GROUP BY "Post".id,"Profile".full_name,"Profile".avatar,"Account".id
		ORDER BY "Post".id DESC    `
		if err := db.Raw(query, account_id).Scan(&result.Posts).Error; err != nil {
			panic(err)
		}
		query = `SELECT COUNT("Post".id) as count_post FROM "Post" WHERE "Post".account_id = ?`
		countPOST := db.Where("account_id =?", account_id).Find(&Post{}).RowsAffected
		db.Raw(`SELECT "Images".* FROM "Post","Images" WHERE "Post".account_id = ? AND "Images".id = ANY("Post".image_ids)`, account_id).Scan(&result.Images)
		response["count_post"] = countPOST
		response["posts"] = result.Posts
		response["images"] = result.Images
		response["is_self"] = true

	} else {
		type Result struct {
			CountPost int      `json:"count_post"`
			Images    []Images `json:"images"`
			FullName  string   `json:"full_name"`
			Avatar    string   `json:"avatar"`
			Posts     []struct {
				ID        uint           `json:"id"`
				Content   string         `json:"content"`
				Avatar    string         `json:"avatar"`
				FullName  string         `json:"full_name"`
				AccountID string         `json:"account_id"`
				UserLiked postgres.Jsonb `json:"user_liked"`
				Comments  postgres.Jsonb `json:"comments"`
				Images    postgres.Jsonb `json:"images"`
			} `json:"posts"`
		}
		var result Result
		response["is_self"] = false
		query := ` SELECT "Post".*,"Profile".avatar,"Profile".full_name , "Account".username,
		(
				SELECT COALESCE(json_agg(DISTINCT item),'[]'::json)
				FROM (
						SELECT "Images".* FROM "Images" WHERE "Images".id = ANY("Post".image_ids)
				) item
		) as images,
		(
				SELECT COALESCE(json_agg(DISTINCT item),'[]'::json)
				FROM ( 
						SELECT "Profile".* FROM "Profile" WHERE "Profile".account_id = ANY("Post".account_liked)) item
		) as user_liked ,
		(
				SELECT COALESCE(json_agg(DISTINCT item),'[]'::json)
				FROM (SELECT "Comment".*,"Profile".* as user FROM "Comment","Profile"
				WHERE "Comment".account_id = "Profile".account_id AND "Comment".post_id = "Post".id
				) as item
		) as comments
		FROM "Post","Profile","Places","Account"
		WHERE "Post".account_id = ? AND "Post".account_id = "Profile".account_id AND "Profile".account_id = "Account".id
		GROUP BY "Post".id,"Profile".full_name,"Profile".avatar,"Account".id
		ORDER BY "Post".id DESC    `
		if err := db.Raw(query, request["account_id"]).Scan(&result.Posts).Error; err != nil {
			panic(err)
		}
		response["posts"] = result.Posts
		query = `SELECT COUNT("Post".id) as count_post FROM "Post" WHERE "Post".account_id = ?`
		countPOST := db.Where("account_id =?", request["account_id"]).Find(&Post{}).RowsAffected
		query = `SELECT "Profile".full_name, "Profile".avatar FROM "Profile" WHERE "Profile".account_id = ?`
		db.Raw(query, request["account_id"]).Scan(&result)
		db.Raw(`SELECT "Images".* FROM "Post","Images" WHERE "Post".account_id = ? AND "Images".id = ANY("Post".image_ids)`, request["account_id"]).Scan(&result.Images)
		response["full_name"] = result.FullName
		response["avatar"] = result.Avatar
		response["images"] = result.Images
		response["count_post"] = countPOST
	}
	// Check friend
	response["is_friend"] = false
	response["status_friend"] = 0
	rows := db.Or("user_1 = ? AND user_2 = ? AND status = 2", request["account_id"], account_id).Or("user_2 = ? AND user_1 = ? AND status = 2", request["account_id"], account_id).Find(&Friends{}).RowsAffected
	if rows > 0 {
		response["is_friend"] = true
		response["status_friend"] = 2
	}
	rows = db.Or("user_1 = ? AND user_2 = ? AND status = 1", request["account_id"], account_id).Or("user_2 = ? AND user_1 = ? AND status = 1", request["account_id"], account_id).Find(&Friends{}).RowsAffected
	if rows > 0 {
		response["is_friend"] = false
		response["status_friend"] = 1
	}

	return c.JSON(http.StatusOK, response)
}

func (*UserController) LoadRequest(c echo.Context) error {
	request := make(map[string]interface{})
	response := echo.Map{}
	account_id, _ := DecodeToken(c, "id")
	if err := c.Bind(&request); err != nil {
		return nil
	}
	var friends []Friends
	var profiles []Profile
	db.Or("user_1 = ? AND status = 1", account_id).Or("user_2 = ? AND status = 1", account_id).Find(&friends)
	if len(friends) == 0 {
		request["profiles"] = []int{}
		return c.JSON(200, request)
	}
	for _, friend := range friends {
		var current Profile

		if account_id.(float64) != float64(friend.User_1) {
			db.Where("account_id = ?", friend.User_1).Find(&current)
			profiles = append(profiles, current)
		} else {
			db.Where("account_id = ?", friend.User_2).Find(&current)
			profiles = append(profiles, current)
		}

	}
	response["profiles"] = profiles
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
	var friend Friends
	if err := db.Or("user_1 = ? AND user_2 = ?", account_id, request["user_id"]).Or("user_1 = ? AND user_2 = ? ", request["user_id"], account_id).Find(&friend).Error; err != nil {
		return c.NoContent(500)
	}
	db.Model(&friend).Update("status", 2)
	return c.JSON(http.StatusOK, response)
}

func (*UserController) CancleAcceptFriend(c echo.Context) error {
	fmt.Println("Here")
	request := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	account_id, _ := DecodeToken(c, "id")
	var friend Friends
	if err := db.Or("user_1 = ? AND user_2 = ?", account_id, request["user_id"]).Or("user_1 = ? AND user_2 = ? ", request["user_id"], account_id).Find(&friend).Error; err != nil {
		return c.NoContent(500)
	}
	friend.Status = 2
	db.Unscoped().Delete(&friend)
	return c.JSON(http.StatusOK, response)
}
func (*UserController) GetUserOnline(c echo.Context) error {
	request := make(map[string]interface{})
	response := make(map[string]interface{})
	if err := c.Bind(&request); err != nil {
		return nil
	}
	account_id, _ := DecodeToken(c, "id")

	type AccountOnline struct {
		FullName  string `json:"full_name"`
		Avatar    string `json:"avatar"`
		AccountID string `json:"account_id"`
		IsOnline  uint   `json:"is_online"`
		SocketID  string `json:"socket_id"`
	}
	var resultForProfile []AccountOnline
	db.Raw(`SELECT "Account".username,"Profile".*,
	CASE
		WHEN "Profile".account_id =  ANY(ARRAY(SELECT "Account_Online".account_id FROM "Account_Online"  WHERE "Account_Online".account_id != ?))
		THEN 1
		ELSE 0
	END AS is_online,
	CASE
		WHEN "Profile".account_id::TEXT = ANY(ARRAY(SELECT "Account_Online".account_id::TEXT FROM "Account_Online")) 
		THEN ((SELECT "Account_Online".socket_id FROM "Account_Online" WHERE "Account_Online".account_id = "Profile".account_id))
		ELSE ''
	END as socket_id
	FROM "Profile","Account"
	WHERE "Account".id = "Profile".account_id AND "Profile".account_id = ANY(ARRAY(SELECT CASE WHEN "Friends".user_1 <> ? THEN "Friends".user_1 ELSE "Friends".user_2 END as current_friend FROM "Friends"
	WHERE ("Friends".user_1 = ? OR "Friends".user_2 =?) AND "Friends".status = 2)) 
	GROUP BY "Profile".id,"Account".id`, account_id, account_id, account_id, account_id).Scan(&resultForProfile)
	response["list_user"] = resultForProfile
	return c.JSON(200, response)
}

func (*UserController) Signout(c echo.Context) error {
	fmt.Println("Logout")
	token, _ := c.Cookie("token")
	token.Value = ""
	token.MaxAge = -1
	c.SetCookie(token)
	return c.String(200, "Logout Success")

}
