// Package sms is the SMS command for the Rabbit service
package sms

import (
	"github.com/aide-family/rabbit/cmd"
	"github.com/spf13/cobra"
)

const cmdLong = `Send SMS messages, supporting multiple SMS service providers and template-based delivery.
The sms command enables direct SMS message sending, supporting configuration of SMS service
providers, recipient phone numbers, message content, and other parameters. It allows quick
single message delivery or template-based sending.

Key Features:
  • SMS delivery: Send SMS messages through configured SMS service providers
  • Multi-provider support: Support for major SMS providers (Alibaba Cloud, Tencent Cloud, Huawei Cloud, etc.)
  • Template-based sending: Support for using provider-provided SMS templates
  • Parameter substitution: Support for dynamic parameter replacement in templates
  • Batch sending: Support for sending SMS messages to multiple phone numbers in batch

Use Cases:
  • SMS testing: Test whether SMS service provider configurations are correct
  • Verification code delivery: Send verification codes, notifications, and other SMS messages
  • Urgent notifications: Send important SMS notifications that require immediate delivery

SMS sending requires prior configuration of SMS service providers (API Key, Secret, etc.),
which can be configured through configuration files or API. Sent SMS messages are processed
immediately, making it suitable for testing and urgent scenarios.`

func NewCmd() *cobra.Command {
	smsCmd := &cobra.Command{
		Use:   "sms",
		Short: "Send SMS messages",
		Long:  cmdLong,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: run,
	}
	smsFlags.addFlags(smsCmd)
	return smsCmd
}

func run(cmd *cobra.Command, args []string) {
	cmd.Help()
}
