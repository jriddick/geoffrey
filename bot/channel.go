package bot

// Channel represents an IRC channel
type Channel struct {
	Name  string
	Users map[string]User
	bot   *Bot
}

// AddUser will add a user to the tracking list
func (c *Channel) AddUser(nick string) {
	if c.Users == nil {
		c.Users = make(map[string]User)
	}

	c.Users[nick] = User{
		Nick: nick,
		bot:  c.bot,
	}
}

// RemoveUser will remove tracked user
func (c *Channel) RemoveUser(nick string) {
	delete(c.Users, nick)
}

// Send will send a message to the channel
func (c *Channel) Send(msg string) {
	c.bot.Send(c.Name, msg)
}
