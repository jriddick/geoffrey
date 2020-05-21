package plugins

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hako/durafmt"
	base "github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
	"github.com/mvdan/xurls"
	log "github.com/sirupsen/logrus"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func init() {
	base.RegisterHandler(YouTubeHandler)
}

const developerKey = ""

// Waiter so we can wait for it to finish before returning
// var wg sync.WaitGroup

// Regex replacer for cleaning titles
// var replacer = regexp.MustCompile("[\r\n]+")

// TitleHandler extracts title from posted links and sends
// them to the channel.
var YouTubeHandler = base.Handler{
	Name:        "YouTube",
	Description: "YouTube parser that extracts title and duration.",
	Event:       irc.Message,
	Run: func(bot *base.Bot, msg *msg.Message) (bool, error) {
		// Get configuration
		config := bot.Config()

		service, err := youtube.NewService(context.Background(), option.WithAPIKey(developerKey))

		if err != nil {
			log.Fatalf("Error creating new YouTube client: %v", err)
		}

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

					URL := uri.String()
					var videoID = strings.Split(URL, "=")[1]
					// Make the API call to YouTube.
					call := service.Videos.List("snippet,contentDetails").Id(videoID)

					response, err := call.Do()

					// Notify on error
					if err != nil {
						log.Errorf("[title] Could not fetch website '%s': %v", uri.String(), err)
					} else {
						for _, video := range response.Items {
							if video.ContentDetails.Duration == "P0D" {
								// fmt.Printf("[YouTube] %v (LIVE)\n", video.Snippet.Title)
								bot.Send(channel, fmt.Sprintf("[%v] %v (%s)", irc.Foreground("YouTube", irc.Green), irc.Bold(video.Snippet.Title), irc.Foreground("LIVE", irc.Orange)))
							} else {
								parsedDur := strings.ToLower(video.ContentDetails.Duration[2:])
								duration, durationErr := durafmt.ParseString(parsedDur)

								if durationErr != nil {
									fmt.Println("Error duration:", durationErr)
								}

								// fmt.Printf("[YouTube] %v (duration: %v)\n", video.Snippet.Title, duration)
								bot.Send(channel, fmt.Sprintf("[%v] %v (duration: %v)", irc.Foreground("YouTube", irc.Green), irc.Bold(video.Snippet.Title), duration))
							}
						}
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
