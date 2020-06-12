package plugins

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	base "github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
	log "github.com/sirupsen/logrus"

	"github.com/tidwall/gjson"
)

func init() {
	base.RegisterHandler(CurrencyHandler)
}

// CurrencyHandler extracts title from posted links and sends
// them to the channel.
var CurrencyHandler = base.Handler{
	Name:        "Currency",
	Description: "Currency converter.",
	Event:       irc.Message,
	Run: func(bot *base.Bot, msg *msg.Message) (bool, error) {
		// Get configuration
		config := bot.Config()

		client := http.Client{}

		returnData := "Currency error, usage: !c <amount> <from> <to>"

		// Check if channel message
		if msg.Params[0] == config.Identification.Nick || msg.Prefix.Name == "nibbler" || msg.Prefix.Name == "geoffrey-bot" {
			return false, nil
		}

		go func(bot *base.Bot, channel string) {
			if strings.HasPrefix(msg.Trailing, "!c") {
				if len(msg.Trailing) > 2 {
					cData := strings.Split(msg.Trailing, " ")

					if len(cData) == 4 {
						url := fmt.Sprintf("https://api.exchangerate.host/latest?base=%s&amount=%s", cData[2], cData[1])

						request, reqErr := http.NewRequest(http.MethodGet, url, nil)
						if reqErr != nil {
							log.Errorf("[currency] Failed NewRequest: %s", reqErr)
						}

						resp, respErr := client.Do(request)
						if respErr != nil {
							log.Errorf("[currency] Failed client.Do: %s", respErr)
						}

						cConvertTo := fmt.Sprintf("rates.%s", strings.ToUpper(cData[3]))

						body, bodyErr := ioutil.ReadAll(resp.Body)
						if bodyErr != nil {
							log.Errorf("[currency] Failed ReadAll: %s", bodyErr)
						}

						cValue := gjson.GetBytes(body, cConvertTo)

						returnData = fmt.Sprintf("%s %s in %s: %v", cData[1], cData[2], cData[3], cValue)
					}
				}
				bot.Send(channel, returnData)
			}
		}(bot, msg.Params[0])
		return true, nil
	},
}
