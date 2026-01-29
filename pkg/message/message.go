// Package message is a simple package that provides a message interface.
package message

import (
	"fmt"
	html "html/template"
	"strings"
	text "text/template"

	"github.com/aide-family/rabbit/pkg/enum"
)

func Equals(m enum.MessageType, channel enum.MessageType) error {
	if m == channel {
		return nil
	}
	return fmt.Errorf("channel %s not supported, only %s is supported", channel, m)
}

type Message interface {
	Marshal() ([]byte, error)
	Type() enum.MessageType
}

func TextFormatter(format string, data any) (string, error) {
	if format == "" {
		return "", fmt.Errorf("format is null")
	}
	if data == nil {
		return "", fmt.Errorf("data is nil")
	}

	t, err := text.New("text/template").Funcs(templateFuncMap).Parse(format)
	if err != nil {
		return "", err
	}
	tpl := text.Must(t, err)
	resultIoWriter := new(strings.Builder)

	if err = tpl.Execute(resultIoWriter, data); err != nil {
		return "", err
	}
	return resultIoWriter.String(), nil
}

func HTMLFormatter(format string, data any) (string, error) {
	if format == "" {
		return "", fmt.Errorf("format is null")
	}
	if data == nil {
		return "", fmt.Errorf("data is nil")
	}

	t, err := html.New("html/template").Funcs(templateFuncMap).Parse(format)
	if err != nil {
		return "", err
	}
	tpl := html.Must(t, err)
	resultIoWriter := new(strings.Builder)

	if err = tpl.Execute(resultIoWriter, data); err != nil {
		return "", err
	}
	return resultIoWriter.String(), nil
}
