package geoffrey

import (
    "github.com/jriddick/geoffrey/irc"
    "github.com/jriddick/geoffrey/msg"
	"fmt"
	"strings"
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
    Channels []string
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

func (b *Bot) Handler() {
    for msg := range b.reader {
        // Send nick and user after connecting
        if msg.Trailing == "*** Looking up your hostname..." {
            b.writer <- fmt.Sprintf("NICK %s", b.config.Nick)
            b.writer <- fmt.Sprintf("USER %s 0 * :%s", b.config.User, b.config.Name)
        }

        // Answer PING with PONG
        if msg.Command == "PING" {
            b.writer <- fmt.Sprintf("PONG %s", msg.Trailing)
        }

        // Join channels when ready
        if msg.Command == irc.RPL_WELCOME {
            for _, channel := range b.config.Channels {
                b.Join(channel)
            }
        }
    }
}

func (b *Bot) Send(channel, msg string) {
    b.writer <- fmt.Sprintf("PRIVMSG %s :%s", channel, msg)
}

func (b *Bot) Join(channel string) {
    // Make sure we have a hashtag
    if !strings.HasPrefix(channel, "#") {
        channel = "#" + channel
    }

    // Send the join command
    b.writer <- fmt.Sprintf("JOIN %s", channel)
}