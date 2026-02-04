// Package gormimpl is the implementation of the gorm repository for the auth service.
package gormimpl

import (
	"context"
	"errors"
	"net/url"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/connect"
	"github.com/aide-family/rabbit/pkg/domain"
	authv1 "github.com/aide-family/rabbit/pkg/domain/auth/v1"
	"github.com/aide-family/rabbit/pkg/domain/auth/v1/gormimpl/model"
	"github.com/aide-family/rabbit/pkg/domain/auth/v1/gormimpl/query"
	"github.com/aide-family/rabbit/pkg/jwt"
	"github.com/aide-family/rabbit/pkg/merr"
)

func init() {
	domain.RegisterAuthV1Factory(config.DomainConfig_GORM, NewGormRepository)
}

func NewGormRepository(c *config.DomainConfig, jwtConfig *config.JWT) (authv1.Repository, func() error, error) {
	ormConfig := &config.ORMConfig{}
	if pointer.IsNotNil(c.GetOptions()) {
		if err := anypb.UnmarshalTo(c.GetOptions(), ormConfig, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, nil, merr.ErrorInternalServer("unmarshal orm config failed: %v", err)
		}
	}
	db, close, err := connect.NewDB(ormConfig)
	if err != nil {
		return nil, nil, err
	}
	query.SetDefault(db)

	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return nil, nil, err
	}
	return &gormRepository{repoConfig: c, db: db, node: node, jwtConfig: jwtConfig}, close, nil
}

type gormRepository struct {
	repoConfig *config.DomainConfig
	db         *gorm.DB
	node       *snowflake.Node
	jwtConfig  *config.JWT
}

// Login implements [authv1.Repository].
func (g *gormRepository) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	user := req.GetUser()
	oauthConfig := req.GetOauthConfig()
	if pointer.IsNil(user) {
		klog.Context(ctx).Debugw("msg", "user is nil")
		return nil, merr.ErrorInvalidArgument("user is nil")
	}
	if pointer.IsNil(oauthConfig) {
		klog.Context(ctx).Debugw("msg", "oauthConfig is nil")
		return nil, merr.ErrorInvalidArgument("oauthConfig is nil")
	}
	if email := user.GetEmail(); strutil.IsEmpty(email) {
		klog.Context(ctx).Debugw("msg", "email is empty")
		return nil, merr.ErrorInvalidArgument("email is empty")
	}

	// 1. check if outh2 user exists
	oauth2UserDO, err := g.findOrCreateOAuth2User(ctx, user)
	if err != nil {
		return nil, err
	}

	// 2. check if user exists
	userDO, err := g.findOrCreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// bind user and oauth2 user
	if err := g.bindUserAndOAuth2User(ctx, userDO, oauth2UserDO); err != nil {
		return nil, err
	}

	// generate token
	token, err := g.generateToken(ctx, userDO)
	if err != nil {
		return nil, err
	}

	// build redirect url
	redirectURL, err := g.buildRedirectURL(token, oauthConfig.GetRedirectURL())
	if err != nil {
		return nil, err
	}
	return &authv1.LoginResponse{
		RedirectURL: redirectURL,
	}, nil
}

func (g *gormRepository) findOrCreateOAuth2User(ctx context.Context, user *authv1.User) (*model.OAuth2User, error) {
	oauth2Mutation := query.OAuth2User
	oauth2UserDO, err := oauth2Mutation.WithContext(ctx).Where(oauth2Mutation.OpenID.Eq(user.GetOpenID())).First()
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Context(ctx).Debugw("msg", "get oauth2 user failed", "error", err, "openID", user.GetOpenID())
			return nil, merr.ErrorInternal("get oauth2 user failed").WithCause(err)
		}
		oauth2UserDO = &model.OAuth2User{
			OpenID: user.GetOpenID(),
			Name:   user.GetName(),
			Email:  user.GetEmail(),
			Avatar: user.GetAvatar(),
			APP:    user.GetApp(),
			Raw:    user.GetRaw(),
			UID:    g.node.Generate(),
		}
		if err := oauth2Mutation.WithContext(ctx).Create(oauth2UserDO); err != nil {
			klog.Context(ctx).Debugw("msg", "create oauth2 user failed", "error", err, "oauth2UserUID", oauth2UserDO.UID)
			return nil, merr.ErrorInternal("create oauth2 user failed").WithCause(err).WithCause(err)
		}
	}
	return oauth2UserDO, nil
}

func (g *gormRepository) findOrCreateUser(ctx context.Context, user *authv1.User) (*model.User, error) {
	userMutation := query.User
	userDO, err := userMutation.WithContext(ctx).Where(userMutation.Email.Eq(user.GetEmail())).First()
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Context(ctx).Debugw("msg", "get user failed", "error", err, "email", user.GetEmail())
			return nil, merr.ErrorInternal("get user failed").WithCause(err)
		}
		userDO = &model.User{
			Email:    user.GetEmail(),
			Name:     user.GetName(),
			Avatar:   user.GetAvatar(),
			Remark:   user.GetRemark(),
			Nickname: user.GetNickname(),
			UID:      g.node.Generate(),
		}
		if err := userMutation.WithContext(ctx).Create(userDO); err != nil {
			klog.Context(ctx).Debugw("msg", "create user failed", "error", err, "userUID", userDO.UID)
			return nil, merr.ErrorInternal("create user failed").WithCause(err)
		}
	}
	return userDO, nil
}

func (g *gormRepository) bindUserAndOAuth2User(ctx context.Context, user *model.User, oauth2User *model.OAuth2User) error {
	if int64(oauth2User.UID) == int64(user.UID) {
		return nil
	}
	oauth2Mutation := query.OAuth2User
	if _, err := oauth2Mutation.WithContext(ctx).Where(oauth2Mutation.UID.Eq(int64(oauth2User.UID))).Update(oauth2Mutation.UID, int64(user.UID)); err != nil {
		klog.Context(ctx).Debugw("msg", "update oauth2 user failed", "error", err, "oauth2UserUID", oauth2User.UID, "userUID", user.UID)
		return merr.ErrorInternal("update oauth2 user failed").WithCause(err)
	}
	return nil
}

func (g *gormRepository) generateToken(ctx context.Context, user *model.User) (string, error) {
	claims := jwt.NewJwtClaims(g.jwtConfig, jwt.BaseInfo{
		UID:      user.UID,
		Username: user.Email,
	})
	return claims.GenerateToken()
}

func (g *gormRepository) buildRedirectURL(token, redirectURL string) (string, error) {
	urlObj, err := url.Parse(redirectURL)
	if err != nil {
		return "", merr.ErrorInvalidArgument("invalid redirect URL").WithCause(err)
	}
	query := urlObj.Query()
	query.Set("token", token)
	urlObj.RawQuery = query.Encode()
	return urlObj.String(), nil
}
