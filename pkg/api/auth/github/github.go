// Package github is the github auth package for the rabbit service.
package github

import (
	"encoding/json"

	"github.com/go-kratos/kratos/v2/transport/http"
	"golang.org/x/oauth2"

	"github.com/aide-family/rabbit/pkg/api/auth"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/merr"
)

func init() {
	auth.RegisterOAuth2LoginFun(config.OAuth2_GITHUB, Login)
}

func Login(ctx http.Context, oauthConfig *oauth2.Config) (auth.User, error) {
	code := ctx.Request().URL.Query().Get("code")
	if code == "" {
		return nil, merr.ErrorInvalidArgument("code is required")
	}
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, merr.ErrorInternal("exchange token failed").WithCause(err)
	}
	client := oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, merr.ErrorInternal("get user info failed").WithCause(err)
	}
	defer resp.Body.Close()
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, merr.ErrorInternal("decode user info failed").WithCause(err)
	}
	return &user, nil
}
