package mockd

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/jriddick/geoffrey/msg"
)

// Mockd is a mocked IRC server
type Mockd struct {
	Port     int
	Listener net.Listener
	Close    chan bool
	sync.WaitGroup
}

// NewMockd takes a port and returns
// a new Mockd object.
func NewMockd(port int) *Mockd {
	return &Mockd{
		Port:  port,
		Close: make(chan bool),
	}
}

// Listen starts the Mockd IRC server
func (m *Mockd) Listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", m.Port))

	if err != nil {
		return err
	}

	m.Listener = listener
	return nil
}

// Handle starts the loop for handling all accepted
// connections.
func (m *Mockd) Handle() {
	m.Add(1)
	defer m.Done()

	conns := m.acceptClientConnections()
	for {
		select {
		case conn := <-conns:
			go m.handleClient(conn)
		case <-m.Close:
			return
		}
	}
}

// Stop will close all connections and stop the handler loop
func (m *Mockd) Stop() (err error) {
	close(m.Close)
	err = m.Listener.Close()
	m.Wait()
	return
}

func (m *Mockd) handleClient(client net.Conn) {
	m.Add(1)
	defer m.Done()

	// Create reader and writer
	reader := bufio.NewReader(client)

	// User registration status
	var userReceived bool
	var nickReceived bool

	// Write opening string
	client.Write([]byte(":geoffrey.com NOTICE Auth :*** Looking up your hostname...\r\n"))

	for {
		select {
		case <-m.Close:
			client.Close()
			return
		default:
			// Read and parse the message
			raw, _ := reader.ReadString('\n')
			msg, _ := msg.ParseMessage(raw)

			if msg != nil {
				// Check if valid NICK command
				if msg.Command == "NICK" {
					nickReceived = true
				}

				// Check if valid USER command
				if msg.Command == "USER" {
					userReceived = true
				}

				if userReceived && nickReceived {
					// Send the registration message
					client.Write([]byte(":geoffrey.com 001 Geoffrey :Welcome to the Geoffrey IRC Network!\r\n"))

					// Unset the status
					nickReceived = false
					userReceived = false
				}

				if msg.Command != "NICK" && msg.Command != "USER" {
					client.Write([]byte(fmt.Sprintf("%s\r\n", raw)))
				}

				continue
			}
		}
	}
}

func (m *Mockd) acceptClientConnections() chan net.Conn {
	conns := make(chan net.Conn)
	go func(conns chan net.Conn, m *Mockd) {
		for {
			client, err := m.Listener.Accept()
			if err != nil || client == nil {
				continue
			}
			conns <- client
		}
	}(conns, m)
	return conns
}
