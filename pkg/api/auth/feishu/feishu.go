// Package feishu is the feishu auth package for the rabbit service.
package feishu

import (
	"encoding/json"
	"errors"

	"github.com/aide-family/magicbox/pointer"
	"github.com/go-kratos/kratos/v2/transport/http"
	"golang.org/x/oauth2"

	"github.com/aide-family/rabbit/pkg/api/auth"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/merr"
)

func init() {
	auth.RegisterOAuth2LoginFun(config.OAuth2_FEISHU, Login)
}

func Login(ctx http.Context, oauthConfig *oauth2.Config) (auth.User, error) {
	code := ctx.Request().URL.Query().Get("code")
	if code == "" {
		return nil, merr.ErrorInvalidArgument("code is required")
	}
	verifier := oauth2.GenerateVerifier()
	token, err := oauthConfig.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, merr.ErrorInternal("exchange token failed").WithCause(err)
	}
	client := oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://open.feishu.cn/open-apis/authen/v1/user_info")
	if err != nil {
		return nil, merr.ErrorInternal("get user info failed").WithCause(err)
	}
	defer resp.Body.Close()
	var userResponse UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, merr.ErrorInternal("decode user info failed").WithCause(err)
	}
	if userResponse.Code != 0 || pointer.IsNil(userResponse.Data) {
		return nil, merr.ErrorInternal("get user info failed").WithCause(errors.New(userResponse.Msg))
	}
	return userResponse.Data, nil
}

type UserResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *User  `json:"data"`
}
