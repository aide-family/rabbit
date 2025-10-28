// Package middler is a package for middleware.
package middler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/magicbox/strutil/cnst"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	jwtv5 "github.com/golang-jwt/jwt/v5"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/merr"
)

func JwtServe(signKey string) middleware.Middleware {
	return jwt.Server(
		func(token *jwtv5.Token) (interface{}, error) {
			return []byte(signKey), nil
		},
		jwt.WithSigningMethod(jwtv5.SigningMethodHS256),
		jwt.WithClaims(func() jwtv5.Claims {
			return &JwtClaims{}
		}),
	)
}

type (
	BaseInfo struct {
		UserID   string `json:"userId"`
		Username string `json:"username"`
	}

	JwtClaims struct {
		signKey string
		BaseInfo
		jwtv5.RegisteredClaims
	}

	baseInfoKey struct{}

	jwtTokenKey struct{}
)

func MustLogin() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			claims, err := GetClaimsFromContext(ctx)
			if err != nil {
				return nil, err
			}
			ctx = WithBaseInfo(ctx, claims.BaseInfo)
			return handler(ctx, req)
		}
	}
}

func BindJwtToken() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			header, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, merr.ErrorUnauthorized("wrong context for middleware")
			}
			auths := strings.SplitN(header.RequestHeader().Get(cnst.HTTPHeaderAuth), " ", 2)
			if len(auths) != 2 || !strings.EqualFold(auths[0], cnst.HTTPHeaderBearerPrefix) {
				return nil, merr.ErrorUnauthorized("token is invalid")
			}

			ctx = WithJwtToken(ctx, auths[1])
			return handler(ctx, req)
		}
	}
}

// GetClaimsFromContext 从context中获取已解析的JWT claims
func GetClaimsFromContext(ctx context.Context) (*JwtClaims, error) {
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return nil, merr.ErrorUnauthorized("token is required")
	}
	jwtClaims, ok := claims.(*JwtClaims)
	if !ok {
		return nil, merr.ErrorUnauthorized("token is invalid")
	}
	return jwtClaims, nil
}

// ParseClaimsFromToken 从JWT token字符串中解析出claims
func ParseClaimsFromToken(secret string, token string) (*JwtClaims, error) {
	claims, err := jwtv5.Parse(token, func(token *jwtv5.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, merr.ErrorUnauthorized("token is invalid")
	}

	claimsBs, err := json.Marshal(claims.Claims)
	if err != nil {
		return nil, err
	}
	var jwtClaims JwtClaims
	if err := json.Unmarshal(claimsBs, &jwtClaims); err != nil {
		return nil, err
	}
	return &jwtClaims, nil
}

func WithBaseInfo(ctx context.Context, baseInfo BaseInfo) context.Context {
	return context.WithValue(ctx, baseInfoKey{}, baseInfo)
}

func GetBaseInfo(ctx context.Context) BaseInfo {
	baseInfo, ok := ctx.Value(baseInfoKey{}).(BaseInfo)
	if !ok {
		return BaseInfo{}
	}
	return baseInfo
}

func WithJwtToken(ctx context.Context, jwtToken string) context.Context {
	return context.WithValue(ctx, jwtTokenKey{}, jwtToken)
}

func GetJwtToken(ctx context.Context) string {
	return ctx.Value(jwtTokenKey{}).(string)
}

// NewJwtClaims new jwt claims
func NewJwtClaims(c *conf.JWT, base BaseInfo) *JwtClaims {
	expire, issuer := c.GetExpire().AsDuration(), c.GetIssuer()
	if expire <= 0 {
		expire = 10 * time.Minute
	}
	if strutil.IsEmpty(issuer) {
		issuer = "rabbit"
	}
	return &JwtClaims{
		signKey:  c.GetSecret(),
		BaseInfo: base,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(expire)),
			Issuer:    issuer,
		},
	}
}

// GenerateToken generate token
func (l *JwtClaims) GenerateToken() (string, error) {
	return jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, l).SignedString([]byte(l.signKey))
}
