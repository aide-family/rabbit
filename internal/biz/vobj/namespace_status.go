// Package vobj is the value object package for the Rabbit service.
package vobj

//go:generate stringer -type=NamespaceStatus -linecomment -output=namespace_status__string.go
type NamespaceStatus int8

const (
	NamespaceStatusUnknown NamespaceStatus = iota
	NamespaceStatusActive
	NamespaceStatusInactive
	NamespaceStatusDeleted
)
