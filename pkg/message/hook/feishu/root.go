package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

	"github.com/aide-family/magicbox/enum"

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

// Signature fills timestamp and sign for Feishu webhook. Always regenerates so timestamp is within 1 hour.
// Algorithm: sign = base64(HMAC-SHA256(key=secret, message=timestamp+"\n"+secret)).
func (m *Message) Signature(secret string) error {
	m.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	payload := m.Timestamp + "\n" + secret
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	m.Sign = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return nil
}
