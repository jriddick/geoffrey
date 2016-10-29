package irc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/jriddick/geoffrey/msg"
)

// IRC client
type IRC struct {
	sync.WaitGroup
	conn         net.Conn
	get          chan *msg.Message
	put          chan string
	end          chan struct{}
	err          chan error
	config       Config
	reconnecting bool
}

// NewIRC returns a new IRC client
func NewIRC(config Config) *IRC {
	return &IRC{
		config: config,
	}
}

func (m *IRC) loopPut() {
	defer m.Done()

	for {
		select {
		case <-m.end:
			return
		case msg, ok := <-m.put:
			// Make sure we received a value
			if !ok {
				m.err <- fmt.Errorf("send channel is invalid")
				return
			}

			if m.conn == nil {
				m.err <- fmt.Errorf("no connection active open")
				return
			}

			// We do not send any empty values
			if msg == "" {
				m.err <- fmt.Errorf("tried to send empty message")
				continue
			}

			// Make sure the suffix is correct
			if !strings.HasSuffix(msg, "\r\n") {
				msg = msg + "\r\n"
			}

			// Set the timeout
			m.conn.SetWriteDeadline(time.Now().Add(time.Second * 30))

			// Send the message to the server
			_, err := m.conn.Write([]byte(msg))

			// Reset the timeout
			m.conn.SetWriteDeadline(time.Time{})

			// Make sure we did not get any errors
			if err != nil {
				m.err <- err
				return
			}
		}
	}
}

func (m *IRC) loopGet() {
	defer m.Done()

	// Reader for the connection
	reader := bufio.NewReaderSize(m.conn, 1024)

	for {
		select {
		case <-m.end:
			return
		default:
			if m.conn != nil {
				// Set the read timeout
				m.conn.SetReadDeadline(time.Now().Add(time.Second * 30))
			}

			// Fetch the message from the server
			raw, err := reader.ReadString('\n')

			if err != nil {
				// Ignore any errors during reconnect
				if !m.reconnecting {
					m.err <- err
				}

				return
			}

			// Reset the timeout
			if m.conn != nil {
				m.conn.SetReadDeadline(time.Time{})
			}

			// Parse the message
			msg, err := msg.ParseMessage(raw)

			if err != nil {
				m.err <- fmt.Errorf("%v [%s]", err, raw)
				return
			}

			// Send the parsed message
			m.get <- msg
		}
	}
}

// Disconnect will disconnect the client
func (m *IRC) Disconnect() {
	// Close the channels
	if m.end != nil {
		close(m.end)
	}

	if m.put != nil {
		close(m.put)
	}

	if m.get != nil {
		close(m.get)
	}

	// Close the connection
	if m.conn != nil {
		m.conn.Close()
	}

	// Reset the connection
	m.conn = nil

	// Wait for loops
	m.Wait()
}

// Connect will connect the client, create new channels if needed
// and start the handler loops.
func (m *IRC) Connect() error {
	// Don't connect if we already are connected
	if m.conn != nil {
		return fmt.Errorf("connection already active")
	}

	// Holds any encountered errors
	var err error

	// Get the hostname
	hostname := m.config.GetHostname()

	if hostname == ":0" || hostname[0] == ':' {
		return fmt.Errorf("need hostname and port to connect")
	}

	// Create the connection
	if m.config.Secure {
		m.conn, err = tls.Dial("tcp", hostname, &tls.Config{
			InsecureSkipVerify: m.config.InsecureSkipVerify,
		})
	} else {
		m.conn, err = net.Dial("tcp", hostname)
	}

	// Check for errors
	if err != nil {
		return err
	}

	// Create the input and output channels
	if m.get == nil {
		m.get = make(chan *msg.Message, 100)
	}

	if m.put == nil {
		m.put = make(chan string, 100)
	}

	if m.end == nil {
		m.end = make(chan struct{})
	}

	if m.err == nil {
		m.err = make(chan error, 100)
	}

	// Start the loops
	m.Add(2)
	go m.loopGet()
	go m.loopPut()

	return nil
}

// Reconnect will disconnect, stop the loops and then call Connect()
func (m *IRC) Reconnect() error {
	// Flag as reconnecting
	m.reconnecting = true

	// Close the connection
	if m.conn != nil {
		m.conn.Close()
	}

	// Reset the connection
	m.conn = nil

	// Close the channel
	close(m.end)

	// Wait until loops complete
	m.Wait()

	// Create the end channel
	m.end = make(chan struct{})

	// Remove the flag
	m.reconnecting = false

	// Connect to the server again
	return m.Connect()
}

// Reader returns channel for reading messages
func (m *IRC) Reader() <-chan *msg.Message {
	return m.get
}

// Writer returns channel for sending messages
func (m *IRC) Writer() chan<- string {
	return m.put
}

// Errors returns channel for reading errors
func (m *IRC) Errors() <-chan error {
	return m.err
}
