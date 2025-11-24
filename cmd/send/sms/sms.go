// Package sms is the SMS command for the Rabbit service
package sms

import (
	"fmt"

	"github.com/aide-family/rabbit/cmd"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sms",
		Short: "Send a SMS message",
		Long: `发送短信消息，支持多种短信服务商和模板发送。

sms 命令用于直接发送短信消息，支持指定短信服务商配置、接收号码、
消息内容等参数，可以快速发送单条短信或使用模板发送。

主要功能：
  • 短信发送：通过配置的短信服务商发送短信消息
  • 多服务商：支持阿里云、腾讯云、华为云等主流短信服务商
  • 模板发送：支持使用服务商提供的短信模板进行发送
  • 参数替换：支持在模板中使用动态参数进行内容替换
  • 批量发送：支持向多个号码批量发送短信

使用场景：
  • 测试短信：测试短信服务商配置是否正确
  • 验证码发送：发送验证码、通知等短信消息
  • 紧急通知：需要立即发送的重要短信通知

短信发送需要先配置短信服务商（API Key、Secret 等），可通过配置
文件或 API 进行配置。发送的短信会立即处理，适合测试和紧急场景。`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: run,
	}
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("sms called")
}
