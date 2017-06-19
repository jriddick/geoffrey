package bot

// Config is the configuration structure for Bot
type Config struct {
	BotName  string `mapstructure:"name"`
	Hostname string `mapstructure:"host"`
	Port     int
	Secure   struct {
		Enable bool
		Verify bool
	}
	Identification struct {
		Nick string
		User string
		Name string
	}
	Channels []string
	Timings  struct {
		Timeout int
	}
	TimeoutLimit int `mapstructure:"retries"`
	Plugins      []string
}
