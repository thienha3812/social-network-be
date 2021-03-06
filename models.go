package main

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

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
	FullName  string `gorm:"type:varchar(255)" json:"full_name"`
	LiveIn    string `json:"live_in"`
	AccountID uint   `gorm:"column:account_id" json:"account_id"`
	Avatar    string `json:"avatar"`
}
type Message struct {
	gorm.Model
	AccountID      uint   `gorm:"not null" json:"account_id"`
	ConversationID uint   `gorm:"not null" json:"conversation_id"`
	Status         uint   `gorm:"default 1" json:"status"`
	Content        string `json:"content"`
}
type Conversation struct {
	gorm.Model
	User_1 uint `gorm:"not null"`
	User_2 uint `gorm:"not null"`
}
type Post struct {
	gorm.Model
	Content      string `gorm:"not null"`
	Like         uint   `gorm:"default:0"`
	AccountID    uint   `gorm:"not null"`
	Comment      []Comment
	ImageIds     pq.Int64Array `gorm:"type:int[]"`
	Rating       float64       `gorm:"default:0"`
	AccountLiked pq.Int64Array `gorm:"type:int[]"`
}
type Comment struct {
	gorm.Model
	Content   string `gorm:"not null"`
	Like      uint   `gorm:"default:0"`
	PostID    uint   `gorm:"not null"`
	AccountID uint   `gorm:"not null"`
}
type Places struct {
	gorm.Model
	Coordinate   string        `gorm:"type:jsonb"`
	Images       pq.Int64Array `gorm:"type:int[]"`
	LandingImage uint
	Description  string
	Name         string
	Address      string
	PostIds      pq.Int64Array `gorm:"type:int[]" sql:"default : '{}'"`
	AccountLiked pq.Int64Array `gorm:"type:int[]"`
}
type Images struct {
	gorm.Model
	Src  string `json:"src"`
	Type string `json:"type"`
	Size uint   `json:"size"`
}
type AccountOnline struct {
	gorm.Model
	AccountID uint `gorm:"not null;column:account_id"`
	IP        string
	SocketID  string  `gorm:"not null"`
	Profile   Profile `gorm:"foreignkey:account_id;association_foreignkey:account_id"`
}

type Reviews struct {
	gorm.Model
	AccountID uint          `gorm:"not null"`
	Rating    float64       `gorm:"default:0"`
	PlaceID   uint          `gorm:"not null"`
	Conent    string        `gorm:"not null"`
	ImagesID  pq.Int64Array `gorm:"type:int[]" sql:"default:'{}'"`
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
func (*Message) TableName() string {
	return "Message"
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

func (*Places) TableName() string {
	return "Places"
}
func (*Images) TableName() string {
	return "Images"
}

func (*Conversation) TableName() string {
	return "Conversation"
}
