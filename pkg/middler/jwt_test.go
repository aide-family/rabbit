package middler_test

import (
	"strings"
	"testing"
	"time"

	"github.com/aide-family/magicbox/strutil/cnst"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/middler"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestGetJwtToken(t *testing.T) {
	claims := middler.NewJwtClaims(&conf.JWT{
		Secret: "xxx",
		Expire: durationpb.New(24 * 365 * time.Hour),
		Issuer: "rabbit-test",
	}, middler.BaseInfo{
		UserID: "test01",
	})
	token, err := claims.GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatalf("GenerateToken returned empty token")
	}
	t.Log(strings.Join([]string{cnst.HTTPHeaderBearerPrefix, token}, " "))
}
