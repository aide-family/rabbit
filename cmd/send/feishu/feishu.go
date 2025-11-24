// Package feishu is the Feishu command for the Rabbit service
package feishu

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

const cmdLong = `Send Feishu messages, supporting text, rich text, and card message formats.
The feishu command enables direct Feishu message sending, supporting configuration of Feishu
Webhook settings, message content, message formats, and other parameters. It allows quick
single message delivery to Feishu groups or users.

Key Features:
  • Message delivery: Send messages to specified groups or users via Feishu Webhook
  • Format support: Support for text, rich text, card, and other message formats
  • Template-based sending: Support for sending messages using pre-configured templates
  • Interactive support: Support for sending interactive message cards (buttons, selectors, etc.)
  • Batch sending: Support for sending messages to multiple groups or users in batch

Use Cases:
  • Message testing: Test whether Feishu Webhook configurations are correct
  • System notifications: Send system alerts, status notifications, etc. to Feishu groups
  • Urgent notifications: Send important Feishu notifications that require immediate delivery

Feishu message sending requires prior configuration of Feishu Webhook URL, which can be
configured through configuration files or API. Sent messages are processed immediately,
making it suitable for testing and urgent scenarios.`

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "feishu",
		Short: "Send Feishu messages",
		Long:  cmdLong,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: run,
	}
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("feishu called")
}
