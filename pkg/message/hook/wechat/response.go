package wechat

import (
	"fmt"
	"io"

	"github.com/aide-family/magicbox/serialize"
)

type response struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (l *response) Error() string {
	if l.ErrCode == 0 {
		return ""
	}
	return fmt.Sprintf("errcode: %d, errmsg: %s", l.ErrCode, l.ErrMsg)
}

func unmarshalResponse(body io.ReadCloser) error {
	var resp response
	if err := serialize.JSONDecoder(body, &resp); err != nil {
		return err
	}
	return &resp
}
