package plugins

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/oauth2"

	"sync"

	"regexp"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v31/github"
	base "github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
	"github.com/mvdan/xurls"
	log "github.com/sirupsen/logrus"
)

func init() {
	base.RegisterHandler(GitHubHandler)
}

// Waiter so we can wait for it to finish before returning
var ghWg sync.WaitGroup

// Holds the authenticated client
var client *github.Client

// Regex replacer for cleaning githubs
var ghReplacer = regexp.MustCompile("[\r\n]+")

// GitHub link matcher
var matcher = regexp.MustCompile("github\\.com/(?P<Username>[a-zA-Z0-9]+)(/(?P<Repository>[a-zA-Z0-9]+))?")

// GitHubHandler extracts information from GitHub
// when a link is posted.
var GitHubHandler = base.Handler{
	Name:        "GitHub",
	Description: "Extracts information from GitHub when a link is posted.",
	Event:       irc.Message,
	Init: func(bot *base.Bot) (bool, error) {
		// Get the configuration
		config := bot.Config()

		// Get the settings
		if settings, ok := config.Settings["github"]; ok {
			if authentication, ok := settings.(map[interface{}]interface{})["authentication"]; ok {
				if key, ok := authentication.(string); ok {
					ctx := context.Background()
					ts := oauth2.StaticTokenSource(
						&oauth2.Token{AccessToken: key},
					)
					tc := oauth2.NewClient(ctx, ts)
					client = github.NewClient(tc)
				}
			}
		}

		if client == nil {
			log.Warnf("[github] Could not get authentication details.")
			client = github.NewClient(nil)
		}

		return true, nil
	},
	Run: func(bot *base.Bot, msg *msg.Message) (bool, error) {
		// Get configuration
		config := bot.Config()

		// Check if channel message
		if msg.Params[0] == config.Identification.Nick {
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

		if client == nil {
			log.Warnf("[github] Could not get authentication details.")
			client = github.NewClient(nil)
		}

		// Open a read transacation to the database
		db.View(func(txn *badger.Txn) error {
			// Download the information from the webpage
			for _, text := range urls {
				match := matcher.FindStringSubmatch(text)
				if len(match) > 0 {
					if _, err := url.Parse(text); err != nil {
						log.Errorf("[github] Could not parse url '%s': %v", text, err)
					} else {
						// Look for the URL in the database
						value, err := txn.Get([]byte(text))

						// Check if it was found or not
						if err != nil {
							if err != badger.ErrKeyNotFound {
								log.Errorf("[github] Could not query the database: %v", err)
							} else {
								log.Infof("[github] Fetching GitHub information for url '%s'", text)

								if match[3] != "" {
									repo, _, err := client.Repositories.Get(context.Background(), match[1], match[3])
									if err != nil {
										log.Errorf("[github] GitHub returned an error response for URL '%s': %v", text, err)
									} else {

										sendMsg := fmt.Sprintf("[%s]", irc.Foreground("GitHub", irc.Green))
										sendMsg += fmt.Sprintf("[%s] %s ", irc.Foreground(repo.GetOrganization().GetLogin(), irc.Blue), repo.GetName())

										commits, _, err := client.Repositories.ListCommits(context.Background(), match[1], match[3], &github.CommitsListOptions{})

										if err != nil {
											log.Errorf("[github] Could not fetch commits for '%s/%s': %v", match[1], match[2], err)
										} else {
											commitMessage := commits[0].GetCommit().GetMessage()
											if len(commitMessage) > 50 {
												commitMessage = commitMessage[0:47]
												commitMessage += "..."
											}
											sendMsg += fmt.Sprintf("(%s) ", irc.Foreground(commitMessage, irc.Orange))
											sendMsg += fmt.Sprintf("(%s) ", irc.Foreground(humanize.Time(repo.GetUpdatedAt().Time), irc.Blue))
											sendMsg += fmt.Sprintf("(%s ⭐) ", irc.Foreground(strconv.Itoa(repo.GetStargazersCount()), irc.Brown))
											sendMsg += fmt.Sprintf("(%s 🍴) ", irc.Foreground(strconv.Itoa(repo.GetForksCount()), irc.Brown))

											bot.Send(msg.Params[0], sendMsg)
										}
									}
								} else {
									user, _, err := client.Users.Get(context.Background(), match[1])

									if err != nil {
										log.Errorf("[github] GitHub returned an error response for URL '%s': %v", text, err)
									} else {
										sendMsg := fmt.Sprintf("[%s]", irc.Foreground("GitHub", irc.Green))
										sendMsg += fmt.Sprintf("[%s] %s ", irc.Foreground(user.GetType(), irc.Blue), user.GetName())
										if user.GetPublicRepos() > 0 {
											sendMsg += fmt.Sprintf("(repositories: %s) ", irc.Foreground(strconv.Itoa(user.GetPublicRepos()), irc.Purple))
										}
										if user.GetPublicGists() > 0 {
											sendMsg += fmt.Sprintf("(gists: %s) ", irc.Foreground(strconv.Itoa(user.GetPublicGists()), irc.Purple))
										}
										if user.GetType() == "User" && user.Company != nil {
											sendMsg += fmt.Sprintf("company: %s)", irc.Foreground(strings.TrimSpace(user.GetCompany()), irc.Teal))
										}
										bot.Send(msg.Params[0], sendMsg)
									}
								}
							}
						} else {
							value.Value(func(val []byte) error {
								bot.Send(msg.Params[0], fmt.Sprintf("[%s] %s",
									irc.Foreground("GitHub", irc.Blue),
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
		ghWg.Wait()

		return true, nil
	},
}
