package bot

import (
	"fmt"
	"strings"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
	"github.com/yuin/gopher-lua"
)

// MessageHandler is the function type for
// message handlers
type MessageHandler func(*Bot, string)

// Bot is the structure for an IRC bot
type Bot struct {
	client          *irc.IRC
	writer          chan<- string
	reader          <-chan *msg.Message
	stop            chan struct{}
	config          Config
	MessageHandlers []MessageHandler
	LuaHandlers     map[string][]*lua.LFunction
	state           *lua.LState
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
	Timeout            int
	TimeoutLimit       int
	ReconnectLimit     int
}

// NewBot creates a new bot
func NewBot(config Config, state *lua.LState) *Bot {
	// Create the bot
	bot := &Bot{
		client: irc.NewIRC(irc.Config{
			Hostname:           config.Hostname,
			Port:               config.Port,
			Secure:             config.Secure,
			InsecureSkipVerify: config.InsecureSkipVerify,
			Timeout:            time.Second * time.Duration(config.Timeout),
			TimeoutLimit:       config.TimeoutLimit,
		}),
		config:      config,
		stop:        make(chan struct{}),
		LuaHandlers: make(map[string][]*lua.LFunction),
		state:       state,
	}

	// Register the bot struct
	RegisterBot(state)
	RegisterConfig(state)

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
	for {
		select {
		case <-b.stop:
			// Send quit message
			b.writer <- "QUIT :Closed"

			// Disconnect the client
			b.client.Disconnect()
			break
		case msg := <-b.reader:
			// Log all messages
			log.Debugln(msg.String())

			// Send nick and user after connecting
			if msg.Trailing == "*** Looking up your hostname..." {
				b.writer <- fmt.Sprintf("NICK %s", b.config.Nick)
				b.writer <- fmt.Sprintf("USER %s 0 * :%s", b.config.User, b.config.Name)
				continue
			}

			// Answer PING with PONG
			if msg.Command == irc.Ping {
				b.writer <- fmt.Sprintf("PONG :%s", msg.Trailing)
				continue
			}

			// Join channels when ready
			if msg.Command == irc.Welcome {
				for _, channel := range b.config.Channels {
					b.Join(channel)
				}
				continue
			}

			// Let our handlers handle PRIVMSG
			if msg.Command == irc.Message {
				// Run the handlers
				go func() {
					for _, handler := range b.MessageHandlers {
						go handler(b, msg.Trailing)
					}
				}()
			}

			// Go through all Lua handlers
			go func(bot *Bot, state *lua.LState) {
				for _, handler := range b.LuaHandlers[msg.Command] {
					// Run the Lua handler
					go func(state *lua.LState, handler *lua.LFunction, bot *Bot) {
						// Push the handler function
						state.Push(handler)

						// Create the metatable for our bot
						value := state.NewUserData()
						value.Value = bot
						state.SetMetatable(value, state.GetTypeMetatable("bot"))

						// Push the bot
						state.Push(value)

						// Push the message
						state.Push(lua.LString(msg.Trailing))

						// Call the function
						state.Call(2, 0)

						// Close the thread
						state.Close()
					}(state.NewThread(), handler, bot)
				}
			}(b, b.state.NewThread())
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
	b.MessageHandlers = append(b.MessageHandlers, handler)
}

// AddLuaHandler registers a Lua handler for the OnMessage event
func (b *Bot) AddLuaHandler(command string, handler *lua.LFunction) {
	b.LuaHandlers[command] = append(b.LuaHandlers[command], handler)
}

// Close will close the bot
func (b *Bot) Close() {
	close(b.stop)
}
