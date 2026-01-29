package auth

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/rabbit/pkg/config"
)

var globalRegistry = NewRegistry()

func NewRegistry() Register {
	return &registry{
		oauth2LoginFuns: safety.NewSyncMap(make(map[config.OAuth2_APP]OAuth2LoginFun)),
	}
}

func RegisterOAuth2LoginFun(app config.OAuth2_APP, loginFun OAuth2LoginFun) {
	globalRegistry.RegisterOAuth2LoginFun(app, loginFun)
}

func GetOAuth2LoginFun(app config.OAuth2_APP) (OAuth2LoginFun, bool) {
	return globalRegistry.GetOAuth2LoginFun(app)
}

type Register interface {
	RegisterOAuth2LoginFun(app config.OAuth2_APP, loginFun OAuth2LoginFun)
	GetOAuth2LoginFun(app config.OAuth2_APP) (OAuth2LoginFun, bool)
}

type registry struct {
	oauth2LoginFuns *safety.SyncMap[config.OAuth2_APP, OAuth2LoginFun]
}

func (r *registry) RegisterOAuth2LoginFun(app config.OAuth2_APP, loginFun OAuth2LoginFun) {
	r.oauth2LoginFuns.Set(app, loginFun)
}

func (r *registry) GetOAuth2LoginFun(app config.OAuth2_APP) (OAuth2LoginFun, bool) {
	return r.oauth2LoginFuns.Get(app)
}
