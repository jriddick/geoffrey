package geoffrey

import (
    "github.com/jriddick/geoffrey/irc"
    "github.com/jriddick/geoffrey/msg"
)

type Bot struct {
    client *irc.IRC
    writer chan<- string
    reader <-chan *msg.Message
    config Config
}

type Config struct {
    irc.Config
    Nick string
    User string
    Name string
}

// Creates a new bot
func NewBot(config Config) *Bot {
    // Create the bot
    bot := &Bot{
        client: irc.NewIRC(config),
        config: config,
    }

    // Get the reader and writer channels
    bot.writer = bot.client.Writer()
    bot.reader = bot.client.Reader()

    return bot
}  
