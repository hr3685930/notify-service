package types

import "github.com/pilagod/gorm-cursor-paginator/v2/cursor"

var UserType = []string{
	"admin",
}

type TestRequest struct {
	ID   int64
	Name string
	Pass string
}

type User struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TestResponse struct {
	Users []*User
	cursor.Cursor
}
