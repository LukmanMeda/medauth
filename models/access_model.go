package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type AccessToken struct {
	models.BaseModel

	Token         string         `db:"token" json:"token"`
	ExpiresAt     types.DateTime `db:"expires_at" json:"expires_at"`
	UserId        string         `db:"user_id" json:"user_id"`
	ClientId      string         `db:"client_id" json:"client_id"`
	IpAddress     string         `db:"ip_address" json:"ip_address"`
	DeviceAddress string         `db:"device_address" json:"device_address"`
	RefreshToken  string         `db:"refresh_token" json:"refresh_token"`
}

type AuthCode struct {
	models.BaseModel

	Code        string `db:"code" json:"code"`
	RedirectUri string `db:"redirect_uri" json:"redirect_uri"`
	UserId      string `db:"user_id" json:"user_id"`
	ClientId    string `db:"client_id" json:"client_id"`
}

type GrantType struct {
	GrantType string `db:"grant_type" json:"grant_type"`
	Note      string `db:"note" json:"note"`
}
