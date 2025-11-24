// Package send is the send command for the Rabbit service
package send

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/cmd/send/email"
	"github.com/aide-family/rabbit/cmd/send/feishu"
	"github.com/aide-family/rabbit/cmd/send/sms"
)

func NewCmd() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send a message to the queue",
		Long: `发送消息到指定通道，支持多种消息类型和发送方式。

send 命令提供直接发送消息的能力，支持邮件、短信、飞书等多种消息
通道，可以快速发送单条消息或使用模板批量发送。

主要功能：
  • 多通道支持：支持邮件、短信、Webhook、飞书等多种消息通道
  • 直接发送：绕过队列直接发送消息，适合紧急或测试场景
  • 模板支持：支持使用预配置的模板进行消息发送
  • 参数验证：发送前自动验证消息参数和通道配置的有效性

子命令：
  • email   发送邮件消息，支持 HTML 和纯文本格式
  • sms     发送短信消息，支持多种短信服务商
  • feishu  发送飞书消息，支持文本和富文本格式

使用场景：
  • 快速测试：测试消息通道配置是否正确
  • 紧急通知：需要立即发送的重要通知
  • 单条发送：发送单条消息，无需通过队列

发送的消息会立即处理，适合测试和紧急场景。对于批量消息发送，
建议使用 apply 命令提交到队列进行异步处理。`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	commands := []*cobra.Command{
		sms.NewCmd(),
		feishu.NewCmd(),
		email.NewCmd(),
	}
	sendCmd.AddCommand(commands...)

	flags.addFlags(sendCmd)

	return sendCmd
}
