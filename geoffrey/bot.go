package geoffrey

import (
	"fmt"
	"strings"

	"log"

	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
)

// MessageHandler is the function type for
// message handlers
type MessageHandler func(*Bot, msg.Message)

// Bot is the structure for an IRC bot
type Bot struct {
	client          *irc.IRC
	writer          chan<- string
	reader          <-chan *msg.Message
	config          Config
	messageHandlers []MessageHandler
}

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
}

// NewBot creates a new bot
func NewBot(config Config) *Bot {
	// Create the bot
	bot := &Bot{
		client: irc.NewIRC(irc.Config{
			Hostname:           config.Hostname,
			Port:               config.Port,
			Secure:             config.Secure,
			InsecureSkipVerify: config.InsecureSkipVerify,
		}),
		config: config,
	}

	return bot
}

// Connect will connect the bot to the server
func (b *Bot) Connect() error {
	// Connect the client
	if err := b.client.Connect(); err != nil {
		return err
	}

	// Get the reader and writer channels
	b.writer = b.client.Writer()
	b.reader = b.client.Reader()

	return nil
}

// Handler will start processing messages
func (b *Bot) Handler() {
	for msg := range b.reader {
		// Log all messages
		log.Println(msg.String())

		// Send nick and user after connecting
		if msg.Trailing == "*** Looking up your hostname..." {
			b.writer <- fmt.Sprintf("NICK %s", b.config.Nick)
			b.writer <- fmt.Sprintf("USER %s 0 * :%s", b.config.User, b.config.Name)
			continue
		}

		// Answer PING with PONG
		if msg.Command == "PING" {
			b.writer <- fmt.Sprintf("PONG %s", msg.Trailing)
			continue
		}

		// Join channels when ready
		if msg.Command == irc.RPL_WELCOME {
			for _, channel := range b.config.Channels {
				b.Join(channel)
			}
			continue
		}

		// Let our handlers handle PRIVMSG
		if msg.Command == "PRIVMSG" {
			go func() {
				for _, handler := range b.messageHandlers {
					handler(b, *msg)
				}
			}()
			continue
		}
	}
}

// Send will send the given message to the given receiver
func (b *Bot) Send(recv, msg string) {
	b.writer <- fmt.Sprintf("PRIVMSG %s :%s", recv, msg)
}

// Join will join the given channel
func (b *Bot) Join(channel string) {
	// Make sure we have a hashtag
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}

	// Send the join command
	b.writer <- fmt.Sprintf("JOIN %s", channel)
}

// OnMessage registeres a new PRIVMSG handler
func (b *Bot) OnMessage(handler MessageHandler) {
	b.messageHandlers = append(b.messageHandlers, handler)
}
