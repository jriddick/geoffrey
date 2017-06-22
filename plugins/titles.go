package plugins

import (
	"fmt"
	"net/url"
	"strings"

	"sync"

	"regexp"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	base "github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
	"github.com/mvdan/xurls"
)

func init() {
	base.RegisterHandler(TitleHandler)
}

// Waiter so we can wait for it to finish before returning
var wg sync.WaitGroup

// Regex replacer for cleaning titles
var replacer = regexp.MustCompile("[\r\n]+")

// TitleHandler extracts title from posted links and sends
// them to the channel.
var TitleHandler = base.Handler{
	Name:        "Title",
	Description: "Extracts title and website information upon detecting URLs",
	Event:       irc.Message,
	Run: func(bot *base.Bot, msg *msg.Message) (bool, error) {
		// Get configuration
		config := bot.Config()

		// Check if channel message
		if msg.Params[0] == config.Identification.Nick || msg.Prefix.Name == "nibbler" {
			return false, nil
		}

		// Extract the urls
		urls := xurls.Relaxed.FindAllString(msg.Trailing, -1)

		// Add the amount of urls needed
		wg.Add(len(urls))

		// Check if we have nothing to do
		if len(urls) < 1 {
			return false, nil
		}

		// Download the information from the webpage
		for _, text := range urls {
			if uri, err := url.Parse(text); err != nil {
				log.Errorf("[title] Could not parse url '%s': %v", text, err)
			} else {
				go func(bot *base.Bot, uri *url.URL, channel string) {
					// Add missing scheme if possible
					if uri.Scheme == "" {
						uri.Scheme = "http"
					}

					// Fetch the document
					doc, err := goquery.NewDocument(uri.String())

					// Notify on error
					if err != nil {
						log.Errorf("[title] Could not fetch website '%s': %v", uri.String(), err)
					} else {
						// Find the title
						bot.Send(channel, fmt.Sprintf("[%s] %s",
							irc.Foreground("LINK", irc.Green),
							irc.Bold(
								strings.TrimSpace(replacer.ReplaceAllString(
									doc.Find("title").Text(),
									" "),
								),
							),
						))
					}

					// Mark as done
					wg.Done()
				}(bot, uri, msg.Params[0])
			}
		}

		// Wait for it to complete
		wg.Wait()

		return true, nil
	},
}
