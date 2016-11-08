package mockd

import (
	"net"
	"testing"

	"bufio"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMockdServer(t *testing.T) {
	Convey("With default Mockd server", t, func() {
		// Create the server
		mockd := NewMockd(5000)

		// Should never be nil
		So(mockd, ShouldNotBeNil)

		// Start accepting connections
		So(mockd.Listen(), ShouldBeNil)

		// Start listening
		go mockd.Handle()

		Convey("It should be able to accept connection", func() {
			// Open connection to mockd
			conn, err := net.Dial("tcp", "127.0.0.1:5000")

			// We should be connected
			So(err, ShouldBeNil)

			// Create buffered reader
			reader := bufio.NewReader(conn)

			// Read entire line
			msg, err := reader.ReadString('\n')

			So(err, ShouldBeNil)
			So(msg, ShouldEqual, ":geoffrey.com NOTICE Auth :*** Looking up your hostname...\r\n")
		})

		Convey("It should be able to handle registration", func() {
			// Open connection to mockd
			conn, err := net.Dial("tcp", "127.0.0.1:5000")

			// We should be connected
			So(err, ShouldBeNil)

			// Create reader
			reader := bufio.NewReader(conn)

			// Read the welcome line
			_, err = reader.ReadString('\n')

			So(err, ShouldBeNil)

			// Send the nick
			_, err = conn.Write([]byte("NICK mockd\r\n"))
			So(err, ShouldBeNil)

			// Send the user
			_, err = conn.Write([]byte("USER mockd * 0 :mockd\r\n"))
			So(err, ShouldBeNil)

			// Read registration proof
			msg, err := reader.ReadString('\n')

			So(err, ShouldBeNil)
			So(msg, ShouldEqual, ":geoffrey.com 001 Geoffrey :Welcome to the Geoffrey IRC Network!\r\n")
		})

		Convey("It should be able to echo anything it receives", func() {
			// Open connection to mockd
			conn, err := net.Dial("tcp", "127.0.0.1:5000")

			// We should be connected
			So(err, ShouldBeNil)

			// Create reader
			reader := bufio.NewReader(conn)

			// Read the welcome line
			_, err = reader.ReadString('\n')

			So(err, ShouldBeNil)

			// Send a message
			_, err = conn.Write([]byte(":geoffrey.com 001 Geoffrey :Welcome to the Geoffrey IRC Network!\r\n"))

			So(err, ShouldBeNil)

			// Read back the same message
			msg, err := reader.ReadString('\n')

			So(err, ShouldBeNil)
			So(msg, ShouldEqual, ":geoffrey.com 001 Geoffrey :Welcome to the Geoffrey IRC Network!\r\n")
		})

		Reset(func() {
			So(mockd.Stop(), ShouldBeNil)
		})
	})
}
