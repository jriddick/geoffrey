package irc

import (
	"fmt"
	"time"
)

// Config is the client configuration
type Config struct {
	Hostname           string
	Port               int
	Secure             bool
	InsecureSkipVerify bool
	Timeout            time.Duration
	TimeoutLimit       int
	MessagesPerSecond  int
}

// GetHostname retuns the full hostname with port
func (c *Config) GetHostname() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}
