package auth

import "github.com/aide-family/rabbit/pkg/config"

type User interface {
	GetOpenID() string
	GetName() string
	GetNickname() string
	GetEmail() string
	GetAvatar() string
	GetAPP() config.OAuth2_APP
	GetRemark() string
	GetRaw() []byte
}
