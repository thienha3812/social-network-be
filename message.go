package main

import (
	fmt "fmt"

	"github.com/labstack/echo/v4"
)

type MessageController struct{}

func (*MessageController) GetHistoryMessage(c echo.Context) error {
	request := echo.Map{}
	response := echo.Map{}
	var friend_message []Message
	var my_message []Message
	fmt.Println(request["friend_id"])
	var conversation Conversation
	if err := c.Bind(&request); err != nil {
		return c.NoContent(500)
	}
	accountID, _ := DecodeToken(c, "id")
	if err := db.Or("user_1 = ? AND user_2 = ?", request["friend_id"], accountID).Or("user_1 = ? AND user_2 = ?", accountID, request["friend_id"]).Find(&conversation).Error; err != nil {
		response["friend_message"] = []string{}
		response["my_message"] = []string{}
		return c.JSON(200, response)
	}
	if err := db.Where("account_id = ? AND conversation_id = ?", accountID, conversation.ID).Find(&my_message).Error; err != nil {
		return c.NoContent(500)
	}
	if err := db.Where("account_id = ? AND conversation_id = ?", request["friend_id"], conversation.ID).Find(&friend_message).Error; err != nil {
		return c.NoContent(500)
	}
	response["friend_message"] = friend_message
	response["my_message"] = my_message
	return c.JSON(200, response)
}
