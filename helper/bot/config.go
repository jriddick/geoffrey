package bot

// Config is the configuration structure for Bot
type Config struct {
	Hostname           string
	Port               int
	Secure             bool
	InsecureSkipVerify bool
	Nick               string
	User               string
	Name               string
	Channels           []string
	Timeout            int
	TimeoutLimit       int
	ReconnectLimit     int
}
