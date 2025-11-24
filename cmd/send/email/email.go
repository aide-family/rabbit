// Package email is the Email command for the Rabbit service
package email

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

func NewCmd() *cobra.Command {
	emailCmd := &cobra.Command{
		Use:   "email",
		Short: "Send an email message",
		Long: `发送邮件消息，支持 HTML 和纯文本格式，可配置邮件服务器。

email 命令用于直接发送邮件消息，支持指定邮件配置 ID、收件人、主题、
内容等参数，可以快速发送单封邮件或使用模板发送。

主要功能：
  • 邮件发送：通过配置的邮件服务器发送邮件消息
  • 格式支持：支持 HTML 和纯文本两种邮件格式
  • 模板发送：支持使用预配置的邮件模板进行发送
  • 附件支持：支持添加邮件附件（如需要）
  • 抄送密送：支持设置收件人、抄送、密送等邮件头

使用场景：
  • 测试邮件：测试邮件服务器配置是否正确
  • 紧急通知：需要立即发送的重要邮件通知
  • 单封发送：发送单封邮件，无需通过队列

邮件发送需要先配置邮件服务器（SMTP），可通过配置文件或 API 进行
配置。发送的邮件会立即处理，适合测试和紧急场景。`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: run,
	}
	flags.addFlags(emailCmd)
	return emailCmd
}
