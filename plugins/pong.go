package plugins

import (
	"strconv"
	"time"

	"github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/msg"
	log "github.com/sirupsen/logrus"
)

func init() {
	bot.RegisterHandler(PongHandler)
}

// PongHandler will handle pong responses
var PongHandler = bot.Handler{
	Name:        "Pong",
	Description: "Handles pong responses from the server",
	Event:       "PONG",
	Run: func(bot *bot.Bot, msg *msg.Message) (bool, error) {
		if num, err := strconv.ParseInt(msg.Trailing, 10, 64); err == nil {
			// Get the time it was sent
			sent := time.Unix(0, num)

			// Log the event
			log.Infof("[pong] Latency is %s", time.Since(sent))
		}

		return true, nil
	},
}
