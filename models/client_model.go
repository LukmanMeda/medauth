package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*Client)(nil)
var _ models.Model = (*ClientUser)(nil)
var _ models.Model = (*Scope)(nil)
var _ models.Model = (*Role)(nil)
var _ models.Model = (*CompanyDepartment)(nil)

type Client struct {
	models.BaseModel

	Name               string         `db:"name" json:"name"`
	CompanyName        string         `db:"company_name" json:"company_name"`
	CompanyWebsite     string         `db:"company_website" json:"company_website"`
	ClientId           string         `db:"client_id" json:"client_id"`
	ClientSecret       string         `db:"client_secret" json:"client_secret"`
	RedirectUri        string         `db:"redirect_uri" json:"redirect_uri"`
	CustomMedaEndpoint string         `db:"custom_meda_endpoint" json:"custom_meda_endpoint"`
	WhiteListDomain    string         `db:"whitelist_domain" json:"whitelist_domain"`
	WhiteListIps       string         `db:"whitelist_ips" json:"whitelist_ips"`
	BilingExpired      types.DateTime `db:"biling_expired" json:"biling_expired"`
	Active             bool           `db:"active" json:"active"`
	LockDevice         bool           `db:"lock_device" json:"lock_device"`
}

type ClientUser struct {
	models.BaseModel

	UserId   string   `db:"user_id" json:"user_id"`
	ClientId string   `db:"client_id" json:"client_id"`
	Scopes   []string `db:"scopes" json:"scopes"`
}

type Scope struct {
	models.BaseModel

	Name        string   `db:"name" json:"name"`
	RoleMinimum int      `db:"role_minimum" json:"role_minimum"`
	Note        string   `db:"note" json:"note"`
	RoleId      []string `db:"role_id" json:"role_id"`
}

type Role struct {
	models.BaseModel

	Name  string `db:"name" json:"name"`
	Value int    `db:"value" json:"value"`
	Note  string `db:"note" json:"note"`
}

type CompanyDepartment struct {
	models.BaseModel

	Title string `db:"title" json:"title"`
}

// TableName implements models.Model
func (*CompanyDepartment) TableName() string {
	return "company_department"
}

// TableName implements models.Model
func (*Client) TableName() string {
	return "clients"
}

// TableName implements models.Model
func (*Role) TableName() string {
	return "roles"
}

// TableName implements models.Model
func (*Scope) TableName() string {
	return "scopes"
}

// TableName implements models.Model
func (*ClientUser) TableName() string {
	return "client_user"
}
