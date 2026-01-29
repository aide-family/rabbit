package feishu

import (
	"fmt"
	"io"

	"github.com/aide-family/magicbox/serialize"
)

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (l *response) Error() string {
	if l.Code == 0 {
		return ""
	}
	return fmt.Sprintf("code: %d, msg: %s, data: %v", l.Code, l.Msg, l.Data)
}

func unmarshalResponse(body io.ReadCloser) error {
	var resp response
	if err := serialize.JSONDecoder(body, &resp); err != nil {
		return err
	}
	return &resp
}
