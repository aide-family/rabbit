package conf

func (c *Bootstrap) IsDev() bool {
	return c.GetEnvironment() == Environment_DEV
}
