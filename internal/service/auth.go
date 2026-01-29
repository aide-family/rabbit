package service

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"golang.org/x/oauth2"

	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/pkg/api/auth"
)

type AuthService struct {
	loginBiz *biz.LoginBiz
}

func NewAuthService(loginBiz *biz.LoginBiz) *AuthService {
	return &AuthService{loginBiz: loginBiz}
}

func (s *AuthService) Login(ctx http.Context, oauthConfig *oauth2.Config, user auth.User) (string, error) {
	return s.loginBiz.Login(ctx, oauthConfig, user)
}
