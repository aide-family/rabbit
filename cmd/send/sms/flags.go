package sms

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd/send"
)

type Flags struct {
	send.SendFlags
}

var smsFlags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.SendFlags = send.GetSendFlags()
}

func GetSmsFlags() Flags {
	return smsFlags
}
