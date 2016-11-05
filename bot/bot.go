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
	client      *irc.IRC
	writer      chan<- string
	reader      <-chan *msg.Message
	stop        chan struct{}
	config      Config
	LuaHandlers map[string][]*lua.LFunction
	state       *lua.LState
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
	RegisterMessage(state)

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

			// Go through all Lua handlers
			go func(bot *Bot, state *lua.LState) {
				for _, handler := range bot.LuaHandlers[msg.Command] {
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
						PushMessage(msg, state)

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

// Ping will send ping to the server
func (b *Bot) Ping(message string) {
	b.writer <- "PING :" + message
}

// Pong will send pong to the server
func (b *Bot) Pong(message string) {
	b.writer <- "PONG :" + message
}

// Nick will send the nick command to the server and
// update the stored nick.
func (b *Bot) Nick(nick string) {
	// Set the nick
	b.config.Nick = nick

	// Send the nick
	b.writer <- "NICK " + nick
}

// User will send the user command to the server and
// update the stored name and user
func (b *Bot) User(user, name string) {
	// Set the stored user and name
	b.config.User = user
	b.config.Name = name

	// Send the command
	b.writer <- "USER " + user + " 0 * :" + name
}

// AddLuaHandler registers a Lua handler for the OnMessage event
func (b *Bot) AddLuaHandler(command string, handler *lua.LFunction) {
	b.LuaHandlers[command] = append(b.LuaHandlers[command], handler)
}

// Close will close the bot
func (b *Bot) Close() {
	close(b.stop)
}
