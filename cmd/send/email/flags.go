package email

import (
	"strings"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/cmd"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/spf13/cobra"
)

type Flags struct {
	cmd.GlobalFlags

	UID         string   `json:"uid" yaml:"uid"`
	Subject     string   `json:"subject" yaml:"subject"`
	Body        string   `json:"body" yaml:"body"`
	To          []string `json:"to" yaml:"to"`
	Cc          []string `json:"cc" yaml:"cc"`
	ContentType string   `json:"contentType" yaml:"contentType"`
	Headers     []string `json:"headers" yaml:"headers"`

	JSON          string `json:"json" yaml:"json"`
	requestParams *apiv1.SendEmailRequest
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	c.Flags().StringVarP(&f.UID, "uid", "u", "", "The uid of the email")
	c.Flags().StringVarP(&f.Subject, "subject", "s", "", "The subject of the email")
	c.Flags().StringVarP(&f.Body, "body", "b", "", "The body of the email")
	c.Flags().StringSliceVarP(&f.To, "to", "t", []string{}, "The to of the email, example: --to=user1@example.com --to=user2@example.com")
	c.Flags().StringSliceVarP(&f.Cc, "cc", "c", []string{}, "The cc of the email, example: --cc=user3@example.com --cc=user4@example.com")
	c.Flags().StringVar(&f.ContentType, "content-type", "text/plain", "The content type of the email")
	c.Flags().StringSliceVarP(&f.Headers, "header", "H", []string{}, "The headers of the email, example: --header=X-Custom-Header:value --header=X-Another-Header:value")
	c.Flags().StringVarP(&f.JSON, "json", "j", "", `{
	"subject": "Test Email",
	"body": "This is a test email",
	"to": ["user1@example.com", "user2@example.com"],
	"cc": ["user3@example.com", "user4@example.com"],
	"contentType": "text/plain",
	"headers": ["X-Custom-Header:value", "X-Another-Header:value"]
}`)
}

func (f *Flags) applyToBootstrap(bc *config.ClientConfig) {
}

func (f *Flags) parseRequestParams() (*apiv1.SendEmailRequest, error) {
	if strutil.IsEmpty(f.JSON) {
		headers := make(map[string]string)
		for _, header := range f.Headers {
			parts := strings.SplitN(header, "=", 2)
			if len(parts) == 2 {
				headers[parts[0]] = parts[1]
			}
		}
		return &apiv1.SendEmailRequest{
			Uid:         f.UID,
			Subject:     f.Subject,
			Body:        f.Body,
			To:          f.To,
			Cc:          f.Cc,
			ContentType: f.ContentType,
			Headers:     headers,
		}, nil
	}
	var requestParams apiv1.SendEmailRequest
	if err := encoding.GetCodec("json").Unmarshal([]byte(f.JSON), &requestParams); err != nil {
		return nil, err
	}
	return &requestParams, nil
}
