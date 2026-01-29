package feishu

import (
	"encoding/json"

	"github.com/aide-family/rabbit/pkg/api/auth"
	"github.com/aide-family/rabbit/pkg/config"
)

var _ auth.User = (*User)(nil)

type User struct {
	Name            string `json:"name"`
	EnName          string `json:"en_name"`
	AvatarURL       string `json:"avatar_url"`
	AvatarThumb     string `json:"avatar_thumb"`
	AvatarMiddle    string `json:"avatar_middle"`
	AvatarBig       string `json:"avatar_big"`
	OpenID          string `json:"open_id"`          // Unique identifier for the user in the current application
	UnionID         string `json:"union_id"`         // Unique identifier for the user in the Feishu open platform
	Email           string `json:"email"`            // User's email
	EnterpriseEmail string `json:"enterprise_email"` // Enterprise email
	ID              string `json:"user_id"`          // User ID (legacy field)
	Mobile          string `json:"mobile"`           // Phone number (with country code)
	TenantKey       string `json:"tenant_key"`       // Enterprise unique identifier
	EmployeeNo      string `json:"employee_no"`      // Employee number
}

// GetAPP implements [auth.User].
func (u *User) GetAPP() config.OAuth2_APP {
	return config.OAuth2_FEISHU
}

// GetAvatar implements [auth.User].
func (u *User) GetAvatar() string {
	return u.AvatarURL
}

// GetEmail implements [auth.User].
func (u *User) GetEmail() string {
	return u.Email
}

// GetName implements [auth.User].
func (u *User) GetName() string {
	return u.Name
}

// GetOpenID implements [auth.User].
func (u *User) GetOpenID() string {
	return u.OpenID
}

// GetRaw implements [auth.User].
func (u *User) GetRaw() []byte {
	raw, _ := json.Marshal(u)
	return raw
}

// GetRemark implements [auth.User].
func (u *User) GetRemark() string {
	return ""
}

// GetNickname implements [auth.User].
func (u *User) GetNickname() string {
	return u.EnName
}
