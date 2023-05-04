package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*User)(nil)
var _ models.Model = (*UserProfile)(nil)

type User struct {
	models.BaseModel
	//models.Record.User

	Username        string `db:"username" json:"username"`
	Email           string `db:"email" json:"email"`
	Name            string `db:"name" json:"name"`
	Avatar          string `db:"avatar" json:"avatar"`
	Phone           string `db:"phone" json:"phone"`
	Password        string `db:"password" json:"password"`
	PasswordConfirm string `db:"passwordConfirm" json:"passwordConfirm"`
	Active          bool   `db:"active" json:"active"`
}

type UserProfile struct {
	models.BaseModel

	Firstname string         `db:"firstname" json:"firstname"`
	Birthday  types.DateTime `db:"birthday" json:"birthday"`
	Lastname  string         `db:"lastname" json:"lastname"`
	UserId    string         `db:"user_id" json:"user_id"`
	Note      string         `db:"note" json:"note"`
}

// TableName implements models.Model
func (*UserProfile) TableName() string {
	return "users_profile"
}

// TableName implements models.Model
func (*User) TableName() string {
	return "users"
}
