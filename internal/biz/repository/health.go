// Package repository is the repository package for the Rabbit service.
package repository

type Health interface {
	Readiness() error
}
