// Package email is the Email command for the Rabbit service
package email

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

const cmdLong = `Send email messages, supporting both HTML and plain text formats with configurable mail servers.

The email command enables direct email message sending, supporting configuration of email
config ID, recipients, subject, content, and other parameters. It allows quick single
email delivery or template-based sending.

Key Features:
  • Email delivery: Send email messages through configured mail servers
  • Format support: Support for both HTML and plain text email formats
  • Template-based sending: Support for sending emails using pre-configured templates
  • Attachment support: Support for adding email attachments (if needed)
  • Email headers: Support for setting recipients, CC, BCC, and other email headers

Use Cases:
  • Email testing: Test whether mail server configurations are correct
  • Urgent notifications: Send important email notifications that require immediate delivery
  • Single email delivery: Send individual emails without going through the queue

Email sending requires prior configuration of mail servers (SMTP), which can be configured
through configuration files or API. Sent emails are processed immediately, making it
suitable for testing and urgent scenarios.`

func NewCmd() *cobra.Command {
	emailCmd := &cobra.Command{
		Use:   "email",
		Short: "Send email messages",
		Long:  cmdLong,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: run,
	}
	flags.addFlags(emailCmd)
	return emailCmd
}
