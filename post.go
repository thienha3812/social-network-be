package main

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type PostController struct{}

func (*PostController) LikePost(c echo.Context) error {
	request := echo.Map{}
	if err := c.Bind(&request); err != nil {
		return c.NoContent(500)
	}
	accountID, _ := DecodeToken(c, "id")
	db.Table("Post").Where("id = ?", request["post_id"]).Update("account_liked", gorm.Expr("array_append(account_liked,?)", accountID))
	return c.NoContent(200)
}

func (*PostController) UserComment(c echo.Context) error {
	request := echo.Map{}
	accountID, _ := DecodeToken(c, "id")
	if err := c.Bind(&request); err != nil {
		return c.NoContent(500)
	}
	post := Comment{
		AccountID: uint(int(accountID.(float64))),
		Content:   request["content"].(string),
		PostID:    uint(int(request["post_id"].(float64))),
	}
	if err := db.Create(&post).Error; err != nil {
		return c.NoContent(500)
	}
	return c.NoContent(200)

}
