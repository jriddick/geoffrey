package irc

import (
	"fmt"
)

type Config struct {
	Hostname           string
	Port               int
	Secure             bool
	InsecureSkipVerify bool
}

func (c *Config) GetHostname() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}
