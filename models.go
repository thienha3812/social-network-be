//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package main


import (
	"time"
)


type ColumnsAccount struct{
	ID, Username, Password, IsOnline, CreatedAt, UpdatedAt, StatusAccount, Email string
}

type ColumnsPost struct{
	ID, AccountID, CreatedAt, UpdatedAt, Content, ImageID, ReplyID, LikedID, Status string
}

type ColumnsProfile struct{
	ID, AccountID, FullName, BirthDay, Logan, Hobbit, LiveIn, Job string
	string
}

type ColumnsSt struct {
	Account ColumnsAccount
	Post ColumnsPost
	Profile ColumnsProfile
}
var Columns = ColumnsSt{
	Account: ColumnsAccount{
		ID: "id",
		Username: "username",
		Password: "password",
		IsOnline: "is_online",
		CreatedAt: "created_at",
		UpdatedAt: "updated_at",
		StatusAccount: "status_account",
		Email: "email",
	},
	Post: ColumnsPost{
		ID: "id",
		AccountID: "account_id",
		CreatedAt: "created_at",
		UpdatedAt: "updated_at",
		Content: "content",
		ImageID: "image_id",
		ReplyID: "reply_id",
		LikedID: "liked_id",
		Status: "status",
	},
	Profile: ColumnsProfile{
		ID: "id",
		AccountID: "account_id",
		FullName: "full_name",
		BirthDay: "birth_day",
		Logan: "logan",
		Hobbit: "hobbit",
		LiveIn: "live_in",
		Job: "job",
	},
}

type TableAccount struct {
	Name, Alias string
}

type TablePost struct {
	Name, Alias string
}

type TableProfile struct {
	Name, Alias string
}

type TablesSt struct {
	Account TableAccount
	Post TablePost
	Profile TableProfile
}
var Tables = TablesSt {
	Account: TableAccount{
		Name: "Account",
		Alias: "t",
	},
	Post: TablePost{
		Name: "Post",
		Alias: "t",
	},
	Profile: TableProfile{
		Name: "Profile",
		Alias: "t",
	},
}

type Account struct {
	tableName struct{} `pg:"\"Account\",alias:t" pg:",discard_unknown_columns"`

	ID int64 `pg:"id,pk"`
	Username string `pg:"username,notnull"`
	Password string `pg:"password,notnull"`
	IsOnline int `pg:"is_online,notnull"`
	CreatedAt time.Time `pg:"created_at,notnull"`
	UpdatedAt time.Time `pg:"updated_at,notnull"`
	StatusAccount int `pg:"status_account,notnull"`
	Email string `pg:"email,notnull"`
}

type Post struct {
	tableName struct{} `pg:"\"Post\",alias:t" pg:",discard_unknown_columns"`

	ID int64 `pg:"id,pk"`
	AccountID int64 `pg:"account_id,notnull"`
	CreatedAt *time.Time `pg:"created_at"`
	UpdatedAt *string `pg:"updated_at"`
	Content *string `pg:"content"`
	ImageID []int `pg:"image_id,array"`
	ReplyID []int `pg:"reply_id,array"`
	LikedID []int `pg:"liked_id,array"`
	Status *int `pg:"status"`
}

type Profile struct {
	tableName struct{} `pg:"\"Profile\",alias:t" pg:",discard_unknown_columns"`

	ID int64 `pg:"id,pk"`
	AccountID int64 `pg:"account_id,notnull"`
	FullName string `pg:"full_name,notnull"`
	BirthDay *string `pg:"birth_day"`
	Logan *string `pg:"logan"`
	Hobbit *string `pg:"hobbit"`
	LiveIn *string `pg:"live_in"`
	Job *string `pg:"job"`

	*Profile `pg:"fk:id"`
}

