// Package conf is the configuration for the Rabbit service.
package conf

func (c *Bootstrap) IsDev() bool {
	return c.GetEnvironment() == Environment_DEV
}
