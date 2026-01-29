package impl

import (
	"context"

	_ "github.com/aide-family/rabbit/pkg/domain/auth/v1/gormimpl"

	klog "github.com/go-kratos/kratos/v2/log"
	"golang.org/x/oauth2"

	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/api/auth"
	"github.com/aide-family/rabbit/pkg/domain"
	authv1 "github.com/aide-family/rabbit/pkg/domain/auth/v1"
	"github.com/aide-family/rabbit/pkg/merr"
)

type loginRepository struct {
	repo authv1.Repository
}

func NewLoginRepository(c *conf.Bootstrap, d *data.Data) (repository.LoginRepository, error) {
	repoConfig := c.GetLoginConfig()
	version := repoConfig.GetVersion()
	driver := repoConfig.GetDriver()
	switch version {
	default:
		factory, ok := domain.GetAuthV1Factory(driver)
		if !ok {
			return nil, merr.ErrorInternalServer("auth repository factory not found")
		}
		repoImpl, close, err := factory(repoConfig, c.GetJwt())
		if err != nil {
			return nil, err
		}
		d.AppendClose("loginRepo", close)

		return &loginRepository{repo: repoImpl}, nil
	}
}

func (l *loginRepository) Login(ctx context.Context, oauthConfig *oauth2.Config, user auth.User) (string, error) {
	req := &authv1.LoginRequest{
		OauthConfig: &authv1.OAuth2Config{
			ClientID:     oauthConfig.ClientID,
			ClientSecret: oauthConfig.ClientSecret,
			RedirectURL:  oauthConfig.RedirectURL,
			Scopes:       oauthConfig.Scopes,
			Endpoint: &authv1.Endpoint{
				AuthURL:       oauthConfig.Endpoint.AuthURL,
				TokenURL:      oauthConfig.Endpoint.TokenURL,
				DeviceAuthURL: oauthConfig.Endpoint.DeviceAuthURL,
				AuthStyle:     int32(oauthConfig.Endpoint.AuthStyle),
			},
		},
		User: &authv1.User{
			OpenID:   user.GetOpenID(),
			Name:     user.GetName(),
			Nickname: user.GetNickname(),
			Email:    user.GetEmail(),
			Avatar:   user.GetAvatar(),
			App:      user.GetAPP().String(),
			Raw:      user.GetRaw(),
			Remark:   user.GetRemark(),
		},
	}
	reply, err := l.repo.Login(ctx, req)
	if err != nil {
		klog.Context(ctx).Debugw("msg", "login failed", "error", err)
		return "", err
	}
	return reply.GetRedirectURL(), nil
}
