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

type Irc struct {
	sync.WaitGroup
	conn         net.Conn
	get          chan *msg.Message
	put          chan string
	end          chan struct{}
	err          chan error
	config       Config
	reconnecting bool
}

func NewIRC(config Config) *Irc {
	return &Irc{
		config: config,
	}
}

func (m *Irc) loop_put() {
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

func (m *Irc) loop_get() {
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

func (m *Irc) Disconnect() {
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

func (m *Irc) Connect() error {
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
	go m.loop_get()
	go m.loop_put()

	return nil
}

func (m *Irc) Reconnect() error {
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

func (m *Irc) Reader() <-chan *msg.Message {
	return m.get
}

func (m *Irc) Writer() chan<- string {
	return m.put
}

func (m *Irc) Errors() <-chan error {
	return m.err
}
