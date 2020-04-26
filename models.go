package main

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"not null"`
	Status   int    `gorm:"default:1"`
}

type Friends struct {
	gorm.Model
	User_1 uint `gorm:"not null"`
	User_2 uint `gorm:"not null"`
	Status uint `gorm:"default:1"`
}
type Profile struct {
	gorm.Model
	FullName  string `gorm:"type:varchar(255)"`
	LiveIn    string
	Account   *Account
	AccountID uint `gorm:"column:account_id"`
	Avatar    string
}
type Post struct {
	gorm.Model
	Content   string `gorm:"not null"`
	Like      uint   `gorm:"default:0"`
	AccountID uint   `gorm:"not null"`
	Comment   []Comment
	Account   *Account
}
type Comment struct {
	gorm.Model
	Content   string `gorm:"not null"`
	Like      uint   `gorm:"default:0"`
	PostID    uint   `gorm:"not null"`
	AccountID uint   `gorm:"not null"`
}

type AccountOnline struct {
	gorm.Model
	AccountID uint `gorm:"not null;column:account_id"`
	IP        string
	SocketID  string  `gorm:"not null"`
	Profile   Profile `gorm:"foreignkey:account_id;association_foreignkey:account_id"`
}

func (*Account) TableName() string {
	return "Account"
}
func (*Friends) TableName() string {
	return "Friends"
}

func (*Profile) TableName() string {
	return "Profile"
}
func (*Comment) TableName() string {
	return "Comment"
}
func (*Post) TableName() string {
	return "Post"
}
func (*AccountOnline) TableName() string {
	return "Account_Online"
}
