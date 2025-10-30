package vobj

//go:generate stringer -type=TemplateApp -linecomment -output=template_app__string.go
type TemplateApp int8

const (
	TemplateAppUnknown         TemplateApp = iota // 未知
	TemplateAppEmail                              // 邮件
	TemplateAppSMS                                // SMS
	TemplateAppWebhookOther                       // Webhook-其他
	TemplateAppWebhookDingTalk                    // Webhook-钉钉
	TemplateAppWebhookWechat                      // Webhook-微信
	TemplateAppWebhookFeishu                      // Webhook-飞书
)

// ToWebhookApp 将 TemplateApp 转换为 WebhookApp（仅适用于 webhook 类型）
func (t TemplateApp) ToWebhookApp() WebhookApp {
	switch t {
	case TemplateAppWebhookOther:
		return WebhookAppOther
	case TemplateAppWebhookDingTalk:
		return WebhookAppDingTalk
	case TemplateAppWebhookWechat:
		return WebhookAppWechat
	case TemplateAppWebhookFeishu:
		return WebhookAppFeishu
	default:
		return WebhookAppUnknown
	}
}

// FromWebhookApp 从 WebhookApp 创建 TemplateApp
func FromWebhookApp(app WebhookApp) TemplateApp {
	switch app {
	case WebhookAppOther:
		return TemplateAppWebhookOther
	case WebhookAppDingTalk:
		return TemplateAppWebhookDingTalk
	case WebhookAppWechat:
		return TemplateAppWebhookWechat
	case WebhookAppFeishu:
		return TemplateAppWebhookFeishu
	default:
		return TemplateAppUnknown
	}
}

// IsWebhookType 判断是否为 webhook 类型
func (t TemplateApp) IsWebhookType() bool {
	return t >= TemplateAppWebhookOther && t <= TemplateAppWebhookFeishu
}

// IsEmailType 判断是否为 email 类型
func (t TemplateApp) IsEmailType() bool {
	return t == TemplateAppEmail
}

// IsSMSType 判断是否为 SMS 类型
func (t TemplateApp) IsSMSType() bool {
	return t == TemplateAppSMS
}
