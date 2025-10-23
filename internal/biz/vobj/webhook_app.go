package vobj

//go:generate stringer -type=WebhookApp -linecomment -output=webhook_app__string.go
type WebhookApp int8

const (
	WebhookAppUnknown  WebhookApp = iota // 未知
	WebhookAppOther                      // 其他
	WebhookAppDingTalk                   // 钉钉
	WebhookAppWechat                     // 微信
	WebhookAppFeishu                     // 飞书
)
