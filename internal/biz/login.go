package biz

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/api/auth"
)

func NewLoginBiz(authRepo repository.LoginRepository) *LoginBiz {
	return &LoginBiz{authRepo: authRepo}
}

type LoginBiz struct {
	authRepo repository.LoginRepository
}

func (b *LoginBiz) Login(ctx context.Context, oauthConfig *oauth2.Config, user auth.User) (string, error) {
	redirectURL, err := b.authRepo.Login(ctx, oauthConfig, user)
	if err != nil {
		return "", err
	}
	return redirectURL, nil
}
