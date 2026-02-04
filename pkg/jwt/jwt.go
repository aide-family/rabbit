// Package jwt is a package for jwt.
package jwt

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"

	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/merr"
)

type (
	// BaseInfo is the base info for the jwt claims.
	BaseInfo struct {
		UID      snowflake.ID `json:"uid"`
		Username string       `json:"username"`
	}

	// JwtClaims is the jwt claims.
	JwtClaims struct {
		signKey string
		BaseInfo
		jwtv5.RegisteredClaims
	}
)

// NewJwtClaims creates a new jwt claims.
func NewJwtClaims(c *config.JWT, base BaseInfo) *JwtClaims {
	expire, issuer := c.GetExpire().AsDuration(), c.GetIssuer()
	if expire <= 0 {
		expire = 10 * time.Minute
	}
	if strutil.IsEmpty(issuer) {
		issuer = "moon"
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

// GenerateToken generates a new jwt token.
func (l *JwtClaims) GenerateToken() (string, error) {
	return jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, l).SignedString([]byte(l.signKey))
}

// GetClaimsFromContext gets the jwt claims from the context.
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

// ParseClaimsFromToken parses the jwt claims from the token string.
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
