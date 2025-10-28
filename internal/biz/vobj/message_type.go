package vobj

//go:generate stringer -type=MessageType -linecomment -output=message_type__string.go
type MessageType int8

const (
	MessageTypeUnknown MessageType = iota // 未知
	MessageTypeEmail                      // 邮件
	MessageTypeWebhook                    // webhook
	MessageTypeSMS                        // SMS
)
