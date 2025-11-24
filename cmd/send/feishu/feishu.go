// Package feishu is the Feishu command for the Rabbit service
package feishu

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "feishu",
		Short: "Send a Feishu message",
		Long: `发送飞书消息，支持文本、富文本和卡片等多种消息格式。

feishu 命令用于直接发送飞书消息，支持指定飞书 Webhook 配置、消息
内容、消息格式等参数，可以快速发送单条消息到飞书群或用户。

主要功能：
  • 消息发送：通过飞书 Webhook 发送消息到指定群或用户
  • 格式支持：支持文本、富文本、卡片等多种消息格式
  • 模板发送：支持使用预配置的消息模板进行发送
  • 交互支持：支持发送可交互的消息卡片（按钮、选择器等）
  • 批量发送：支持向多个群或用户批量发送消息

使用场景：
  • 测试消息：测试飞书 Webhook 配置是否正确
  • 系统通知：发送系统告警、状态通知等到飞书群
  • 紧急通知：需要立即发送的重要飞书通知

飞书消息发送需要先配置飞书 Webhook URL，可通过配置文件或 API
进行配置。发送的消息会立即处理，适合测试和紧急场景。`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: run,
	}
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("feishu called")
}
