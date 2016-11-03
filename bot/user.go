package bot

// User represents an IRC user
type User struct {
	Nick string
	bot  *Bot
}

// Send will send a message to the user
func (u *User) Send(msg string) {
	u.bot.Send(u.Nick, msg)
}
