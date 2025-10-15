package email

import (
	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/spf13/cobra"
)

type Flags struct {
	cmd.GlobalFlags

	configPath string

	Subject     string   `json:"subject" yaml:"subject"`
	Body        string   `json:"body" yaml:"body"`
	To          []string `json:"to" yaml:"to"`
	Cc          []string `json:"cc" yaml:"cc"`
	ContentType string   `json:"contentType" yaml:"contentType"`
	Headers     []string `json:"headers" yaml:"headers"`
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.GlobalFlags = cmd.GetGlobalFlags()
	c.Flags().StringVar(&f.configPath, "config", "~/.rabbit/config", "The config of the email")
	c.Flags().StringVarP(&f.Subject, "subject", "s", "", "The subject of the email")
	c.Flags().StringVarP(&f.Body, "body", "b", "", "The body of the email")
	c.Flags().StringSliceVarP(&f.To, "to", "t", []string{}, "The to of the email, example: --to=user1@example.com --to=user2@example.com")
	c.Flags().StringSliceVarP(&f.Cc, "cc", "c", []string{}, "The cc of the email, example: --cc=user3@example.com --cc=user4@example.com")
	c.Flags().StringVar(&f.ContentType, "content-type", "text/plain", "The content type of the email")
	c.Flags().StringSliceVarP(&f.Headers, "header", "H", []string{}, "The headers of the email, example: --header=X-Custom-Header:value --header=X-Another-Header:value")
}

func (f *Flags) applyToBootstrap(bc *conf.Bootstrap) {
}
