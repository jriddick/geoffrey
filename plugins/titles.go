package plugins

import (
	"fmt"
	"net/url"
	"strings"

	"sync"

	"regexp"

	"github.com/PuerkitoBio/goquery"
	badger "github.com/dgraph-io/badger/v2"
	base "github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
	"github.com/mvdan/xurls"
	log "github.com/sirupsen/logrus"
)

func init() {
	base.RegisterHandler(TitleHandler)
}

func fetchTitle(bot *base.Bot, uri *url.URL, channel string, text string) {
	// Add missing scheme if possible
	if uri.Scheme == "" {
		uri.Scheme = "http"
	}

	// Get the database from the bot
	db := bot.Db()

	// Fetch the document
	doc, err := goquery.NewDocument(uri.String())

	// Notify on error
	if err != nil {
		log.Errorf("[title] Could not fetch website '%s': %v", uri.String(), err)
	} else {
		// Clean the title
		title := strings.TrimSpace(replacer.ReplaceAllString(
			doc.Find("title").First().Text(),
			" "),
		)

		// Find the title
		bot.Send(channel, fmt.Sprintf("[%s] %s",
			irc.Foreground("LINK", irc.Green),
			irc.Bold(
				title,
			),
		))

		// Save the title for future use
		if err := db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(text), []byte(title))
		}); err != nil {
			log.Errorf("[title] Could not save title to database: %v", err)
		}
	}

	// Mark as done
	wg.Done()
}

// Waiter so we can wait for it to finish before returning
var wg sync.WaitGroup

// Blacklist for links that should not be handled
var blacklist []*regexp.Regexp

// Regex replacer for cleaning titles
var replacer = regexp.MustCompile("[\r\n]+")

// TitleHandler extracts title from posted links and sends
// them to the channel.
var TitleHandler = base.Handler{
	Name:        "Title",
	Description: "Extracts title and website information upon detecting URLs",
	Event:       irc.Message,
	Init: func(bot *base.Bot) (bool, error) {
		// Get the configuration
		config := bot.Config()

		// Get the blacklist list from the configs
		if settings, ok := config.Settings["title"]; ok {
			if list, ok := settings.(map[interface{}]interface{})["blacklist"]; ok {
				for _, matcher := range list.([]interface{}) {
					if regex, err := regexp.Compile(matcher.(string)); err != nil {
						log.Errorf("[title] Could not compile regex '%s'", matcher)
					} else {
						blacklist = append(blacklist, regex)
					}
				}
			}
		}

		return true, nil
	},
	Run: func(bot *base.Bot, msg *msg.Message) (bool, error) {
		// Get configuration
		config := bot.Config()

		// Check if channel message
		if msg.Params[0] == config.Identification.Nick || msg.Prefix.Name == "nibbler" || msg.Prefix.Name == "geoffrey-bot" {
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

		// Get the database
		db := bot.Db()

		// Open a read transacation to the database
		db.View(func(txn *badger.Txn) error {
			// Download the information from the webpage
			for _, text := range urls {
				// Set to true to to skip the urls from being handled
				skip := false

				// Go through all blacklists
				for _, matcher := range blacklist {
					if matcher.MatchString(text) {
						log.Infof("[title] Skipped link '%s' due to blacklist.", text)
						skip = true
						break
					}
				}

				if !skip {
					if uri, err := url.Parse(text); err != nil {
						log.Errorf("[title] Could not parse url '%s': %v", text, err)
					} else {
						// Look for the URL in the database
						value, err := txn.Get([]byte(text))

						// Check if it was found or not
						if err != nil {
							if err != badger.ErrKeyNotFound {
								log.Errorf("[title] Could not query the database: %v", err)
							} else {
								log.Infof("[title] Fetching title for url '%s'", text)

								// Fetch the title from the website
								go fetchTitle(bot, uri, msg.Params[0], text)
							}
						} else {
							value.Value(func(val []byte) error {
								bot.Send(msg.Params[0], fmt.Sprintf("[%s] %s",
									irc.Foreground("LINK", irc.Green),
									irc.Bold(
										string(val),
									),
								))
								return nil
							})
						}
					}
				}
			}

			return nil
		})

		// Wait for it to complete
		wg.Wait()

		return true, nil
	},
}
