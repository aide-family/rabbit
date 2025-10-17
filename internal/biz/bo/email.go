// Package bo is the business logic object
package bo

type SendEmailBo struct {
	Namespace   string
	Subject     string
	Body        string
	To          []string
	Cc          []string
	ContentType string
	Headers     []string
}
