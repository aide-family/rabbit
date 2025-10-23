// Package vobj is the value object package for the Rabbit service.
package vobj

//go:generate stringer -type=GlobalStatus -linecomment -output=global_status__string.go
type GlobalStatus int8

const (
	GlobalStatusUnknown  GlobalStatus = iota // 未知
	GlobalStatusEnabled                      // 启用
	GlobalStatusDisabled                     // 禁用
)
