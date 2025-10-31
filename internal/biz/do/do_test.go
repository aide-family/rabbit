package do_test

import (
	"testing"
	"time"

	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
)

func TestBaseModel_UID(t *testing.T) {
	nodeID := strutil.GetNodeIDFromIP()
	t.Log(nodeID)
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		t.Fatalf("failed to create snowflake node: %v", err)
	}
	uid := node.Generate()
	int64Val := uid.Int64()
	t.Log(int64Val)
	t.Log(uid.Time())
	t.Log(uid.Node())
	sec := uid.Time()
	t.Log(sec)
	ts := time.UnixMilli(sec)
	t.Log(ts)
	base64 := uid.Base64()
	t.Log(base64)

	t.Log(snowflake.ParseInt64(int64Val))
}
