package irc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"sync"

	"log"

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
		get:    make(chan *msg.Message),
		put:    make(chan string, 100),
		end:    make(chan struct{}),
		err:    make(chan error, 100),
	}
}

func (m *IRC) loopPut() {
	defer m.Done()

	for {
		select {
		case <-m.end:
			return
		case msg := <-m.put:
			log.Println("Sending: ", msg)

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
			//m.conn.SetWriteDeadline(time.Now().Add(m.config.Timeout))

			// Send the message to the server
			_, err := m.conn.Write([]byte(msg))

			// Reset the timeout
			//m.conn.SetWriteDeadline(time.Time{})

			log.Println("Sent: ", msg, "Error: ", err)

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

	// The current amount of timeouts
	timeouts := 0

	for {
		select {
		case <-m.end:
			return
		default:
			// Set the read timeout
			//m.conn.SetReadDeadline(time.Now().Add(m.config.Timeout))

			// Fetch the message from the server
			raw, err := reader.ReadString('\n')

			if err != nil {
				if !m.reconnecting {
					// We continue on timeout
					if err, ok := err.(net.Error); ok && err.Timeout() {
						// Increase the amount of timeouts
						timeouts++

						// Only continue if we haven't reached the limit
						if timeouts < m.config.TimeoutLimit {
							continue
						}
					}

					// Send the error
					m.err <- err
				}

				return
			}

			// Reset the timeout
			//m.conn.SetReadDeadline(time.Time{})

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

	// Wait for loops
	m.Wait()

	// Close the connection
	if m.conn != nil {
		m.conn.Close()
	}

	// Reset the connection
	m.conn = nil

	// Close the put channel
	if m.put != nil {
		close(m.put)
	}

	// Close the get channel
	if m.get != nil {
		close(m.get)
	}
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
