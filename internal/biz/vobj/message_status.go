package vobj

//go:generate stringer -type=MessageStatus -linecomment -output=message_status__string.go
type MessageStatus int8

const (
	MessageStatusUnknown   MessageStatus = iota // 未知
	MessageStatusPending                        // 待处理
	MessageStatusSending                        // 发送中
	MessageStatusSent                           // 已发送
	MessageStatusFailed                         // 失败
	MessageStatusCancelled                      // 已取消
)
