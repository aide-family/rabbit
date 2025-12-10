// Package send is the send command for the Rabbit service
package send

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

const cmdLong = `Send messages to specified channels, supporting multiple message types and delivery methods.

The send command provides direct message sending capabilities, supporting multiple message
channels such as email, SMS, Feishu, etc. It enables quick single message delivery or
batch sending using templates.

Key Features:
  • Multi-channel support: Support for email, SMS, Webhook, Feishu, and other message channels
  • Direct delivery: Bypass the queue for immediate message sending, suitable for urgent or testing scenarios
  • Template support: Support for sending messages using pre-configured templates
  • Parameter validation: Automatically validate message parameters and channel configuration validity before sending

Subcommands:
  • email   Send email messages, supporting both HTML and plain text formats
  • sms     Send SMS messages, supporting multiple SMS service providers
  • feishu  Send Feishu messages, supporting text and rich text formats

Use Cases:
  • Quick testing: Test whether message channel configurations are correct
  • Urgent notifications: Send important notifications that require immediate delivery
  • Single message delivery: Send individual messages without going through the queue

Messages sent through this command are processed immediately, making it suitable for
testing and urgent scenarios. For bulk message sending, it is recommended to use the
apply command to submit messages to the queue for asynchronous processing.`

func NewCmd(children ...*cobra.Command) *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send messages to specified channels",
		Long:  cmdLong,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	sendFlags.addFlags(sendCmd)
	sendCmd.AddCommand(children...)

	return sendCmd
}
