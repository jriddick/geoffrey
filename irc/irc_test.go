package irc

import (
	"testing"
	"time"

	"github.com/jriddick/geoffrey/mockd"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	defaultConfig = Config{
		Hostname:           "127.0.0.1",
		Port:               5000,
		Secure:             false,
		InsecureSkipVerify: false,
		Timeout:            time.Second * 30,
		TimeoutLimit:       5,
		MessagesPerSecond:  120,
	}
)

func TestClient(t *testing.T) {
	// Create the mockd server
	mockd := mockd.NewMockd(5000)

	// Start listening
	mockd.Listen()

	// Start accepting connections
	go mockd.Handle()

	Convey("With the default IRC client", t, func() {
		Convey("It should be able to connect", func() {
			// Open client
			client := NewIRC(defaultConfig)

			// Connect the client
			So(client.Connect(), ShouldBeNil)

			// Get the reader channel
			reader := client.Reader()

			// Read from the server
			So(<-reader, ShouldNotBeNil)
		})

		Convey("It should be able to register", func() {
			// Open client
			client := NewIRC(defaultConfig)

			// Connect to the server
			So(client.Connect(), ShouldBeNil)

			// Get the reader and writer channel
			reader := client.Reader()
			writer := client.Writer()

			// Read the first message
			So(<-reader, ShouldNotBeNil)

			// Send our registration
			writer <- "NICK geoffrey"
			writer <- "USER geoffrey 0 * :geoffrey"

			// Wait for our verification
			result := <-reader

			// Verify that we got what we needed
			So(result, ShouldNotBeNil)
			So(result.String(), ShouldEqual, ":geoffrey.com 001 Geoffrey :Welcome to the Geoffrey IRC Network!")
		})

		Convey("It should be able to reconnect", func() {
			// Open client
			client := NewIRC(defaultConfig)

			// Connect
			So(client.Connect(), ShouldBeNil)

			// Get the reader
			reader := client.Reader()

			// Wait until we get opening
			So(<-reader, ShouldNotBeNil)

			// Reconnect
			So(client.Reconnect(), ShouldBeNil)

			// Wait until we get opening
			So(<-reader, ShouldNotBeNil)
		})

		Convey("It should reject empty hostname", func() {
			// Open client
			client := NewIRC(Config{
				Hostname:           "",
				Port:               defaultConfig.Port,
				Secure:             defaultConfig.Secure,
				InsecureSkipVerify: defaultConfig.InsecureSkipVerify,
				Timeout:            defaultConfig.Timeout,
				TimeoutLimit:       defaultConfig.TimeoutLimit,
				MessagesPerSecond:  defaultConfig.MessagesPerSecond,
			})

			// Should fail to connect
			So(client.Connect(), ShouldNotBeNil)
		})

		Convey("It should fail to connect securely to non-ssl server", func() {
			// Create the client
			client := NewIRC(Config{
				Hostname:           defaultConfig.Hostname,
				Port:               defaultConfig.Port,
				Secure:             true,
				InsecureSkipVerify: defaultConfig.InsecureSkipVerify,
				Timeout:            defaultConfig.Timeout,
				TimeoutLimit:       defaultConfig.TimeoutLimit,
				MessagesPerSecond:  defaultConfig.MessagesPerSecond,
			})

			// Should fail to connect
			So(client.Connect(), ShouldNotBeNil)
		})

		Convey("It should not be able to connect twice", func() {
			// Create the client
			client := NewIRC(defaultConfig)

			// Should succeed to connect
			So(client.Connect(), ShouldBeNil)

			// Should not succeed a second time
			So(client.Connect(), ShouldNotBeNil)
		})

		Convey("It should not be able to send empty messages", func() {
			// Create the client
			client := NewIRC(defaultConfig)

			// Should succeed to connect
			So(client.Connect(), ShouldBeNil)

			// Get the writer
			writer := client.Writer()

			// Send empty message
			writer <- ""

			// Get error channel
			errors := client.Errors()

			// We should receive error
			err := <-errors

			// Check the error
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "tried to send empty message")
		})

		Convey("It should be able to disconnect", func() {
			// Create the client
			client := NewIRC(Config{
				Hostname:           defaultConfig.Hostname,
				Port:               defaultConfig.Port,
				Secure:             defaultConfig.Secure,
				InsecureSkipVerify: defaultConfig.InsecureSkipVerify,
				Timeout:            time.Second * 1,
				TimeoutLimit:       defaultConfig.TimeoutLimit,
				MessagesPerSecond:  defaultConfig.MessagesPerSecond,
			})

			// Should connect
			So(client.Connect(), ShouldBeNil)

			// Should disconnect
			client.Disconnect("Leaving")
		})
	})
}
