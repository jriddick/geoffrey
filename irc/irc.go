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
	connected    bool
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

	// Calculate the duration we have to wait to honor MessagesPerSecond
	wait := time.Duration(1000/m.config.MessagesPerSecond) * time.Millisecond

	for {
		select {
		case <-m.end:
			return
		case msg := <-m.put:
			// Wait for the pre-determined time before sending
			time.Sleep(wait)

			// We do not send any empty values
			if msg == "" {
				m.err <- fmt.Errorf("[geoffrey] Tried to send empty message")
				continue
			}

			// Make sure the suffix is correct
			if !strings.HasSuffix(msg, "\r\n") {
				msg = msg + "\r\n"
			}

			// Set the timeout
			m.conn.SetWriteDeadline(time.Now().Add(m.config.Timeout))

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
			// Set the read timeout
			m.conn.SetReadDeadline(time.Now().Add(m.config.Timeout))

			// Fetch the message from the server
			raw, err := reader.ReadString('\n')

			// Make sure we don't have any reading errors
			if err != nil {
				// Send the error that occured if we aren't currently in the process of connecting
				if !m.reconnecting {
					m.err <- err
				}

				return
			}

			// Reset the timeout
			m.conn.SetReadDeadline(time.Time{})

			// Parse the message
			msg, err := msg.ParseMessage(raw)

			if err != nil {
				m.err <- fmt.Errorf("[parse] Could not parse '%s': %v", raw, err)
				return
			}

			// Send the parsed message
			m.get <- msg
		}
	}
}

// Disconnect will disconnect the client
func (m *IRC) Disconnect(message string) {
	// Close the channels
	if m.end != nil {
		close(m.end)
	}

	// Wait for loops
	m.Wait()

	// Send quit message to the connection
	m.conn.Write([]byte("QUIT :" + message + "\r\n"))

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
		return fmt.Errorf("[geoffrey] Connection already active")
	}

	// Holds any encountered errors
	var err error

	// Get the hostname
	hostname := m.config.GetHostname()

	if hostname == ":0" || hostname[0] == ':' {
		return fmt.Errorf("[geoffrey] Need hostname and port to connect")
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
	// Only reconnect if we are connected
	if m.reconnecting {
		return nil
	}

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
