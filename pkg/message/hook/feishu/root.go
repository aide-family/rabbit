package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

	"github.com/aide-family/rabbit/pkg/enum"
	"github.com/aide-family/rabbit/pkg/message"
)

var _ message.Message = (*Message)(nil)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypePost  MessageType = "post"
	MessageTypeImage MessageType = "image"
	MessageTypeCard  MessageType = "interactive"
)

type Content struct {
	Text  *Text  `json:"text,omitempty"`
	Post  *Post  `json:"post,omitempty"`
	Card  *Card  `json:"card,omitempty"`
	Image string `json:"image_key,omitempty"`
}

type Message struct {
	MsgType   MessageType `json:"msg_type"`
	Content   *Content    `json:"content,omitempty"`
	Timestamp string      `json:"timestamp"`
	Sign      string      `json:"sign"`
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) Type() enum.MessageType {
	return enum.MessageType_WEBHOOK_FEISHU
}

func (m *Message) Signature(secret string) error {
	if m.Timestamp != "" && m.Sign != "" {
		return nil
	}
	m.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	// timestamp + key sha256, then base64 encode
	signString := m.Timestamp + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(signString))
	_, err := h.Write(data)
	if err != nil {
		return err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	m.Sign = signature
	return nil
}
