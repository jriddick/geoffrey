package irc_test

import (
	"net"

	"log"

	"bufio"

	. "github.com/jriddick/geoffrey/irc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func IRCDaemonConn(conn net.Listener, result chan bool) {
	// Accepted the client
	_, err := conn.Accept()
	if err != nil {
		result <- false
	}
	result <- true
}

func IRCDaemonReg(conn net.Listener, result chan []byte) {
	// Accept the client
	client, err := conn.Accept()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connection Established")

	// Write some data
	client.Write([]byte(":geoffrey.com NOTICE Auth :*** Looking up your hostname...\r\n"))

	go func() {
		// Wait until we get some data
		reader := bufio.NewReader(client)
		msg, _, _ := reader.ReadLine()

		log.Println("Message Read: ", string(msg))

		// Return RPL_WELCOME
		client.Write([]byte(":geoffrey.com 001 Geoffrey :Welcome to the Geoffrey IRC Network!\r\n"))

		// Send what we got on the channel
		result <- msg
	}()
}

func IRCDaemonRec(conn net.Listener) {
	// Accept the client
	client, err := conn.Accept()
	if err != nil {
		log.Fatalln(err)
	}

	// Write some data
	client.Write([]byte(":geoffrey.com NOTICE Auth :*** Looking up your hostname...\r\n"))

	// Accept the client
	client2, err := conn.Accept()
	if err != nil {
		log.Fatalln(err)
	}

	// Write some data
	client2.Write([]byte(":geoffrey.com NOTICE Auth :*** Looking up your hostname...\r\n"))
}

var _ = Describe("Irc", func() {
	var (
		cfg Config = Config{
			Hostname:           "127.0.0.1",
			Port:               5555,
			Secure:             false,
			InsecureSkipVerify: false,
			Timeout:            30,
			TimeoutLimit:       5,
		}
		emp Config = Config{}
	)

	It("should be able to create an irc client", func() {
		Expect(NewIRC(cfg)).NotTo(BeNil())
	})

	It("should be able to connect", func() {
		By("starting the server")
		ircd, err := net.Listen("tcp", "127.0.0.1:5555")
		Expect(err).NotTo(HaveOccurred())

		By("starting the handler")
		result := make(chan bool)
		go IRCDaemonConn(ircd, result)

		// Connect the client
		client := NewIRC(cfg)
		Expect(client.Connect()).NotTo(HaveOccurred())
		client.Disconnect()

		// Get the connection results
		connected := <-result
		Expect(connected).To(BeTrue())

		// Close the server
		Expect(ircd.Close()).NotTo(HaveOccurred())
	})

	It("should not be able to connect with an empty configuration", func() {
		Expect(NewIRC(emp).Connect()).To(HaveOccurred())
	})

	It("should be able to register to the server", func() {
		By("creating the server")
		ircd, err := net.Listen("tcp", "127.0.0.1:5556")
		Expect(err).NotTo(HaveOccurred())

		By("starting the handler")
		result := make(chan []byte)
		go IRCDaemonReg(ircd, result)

		By("connecting to the server")
		cfg.Port = 5556
		client := NewIRC(cfg)
		Expect(client.Connect()).NotTo(HaveOccurred())

		By("getting the reader and writer")
		reader := client.Reader()
		Expect(reader).NotTo(BeNil())
		writer := client.Writer()
		Expect(writer).NotTo(BeNil())

		By("sending user")
		writer <- "USER geoffrey 0 * :geoffrey"

		By("checking sent data")
		res := <-result
		Expect(string(res)).To(Equal("USER geoffrey 0 * :geoffrey"))

		Expect(ircd.Close()).NotTo(HaveOccurred())
	})

	It("should be able to reconnect to the server", func() {
		By("creating the server")
		ircd, err := net.Listen("tcp", "127.0.0.1:5557")
		Expect(err).NotTo(HaveOccurred())

		By("starting the handler")
		go IRCDaemonRec(ircd)

		By("connecting to the server")
		cfg.Port = 5557
		client := NewIRC(cfg)
		Expect(client.Connect()).NotTo(HaveOccurred())

		By("reading from the server")
		Expect(<-client.Reader()).NotTo(BeNil())

		By("reconnecting")
		Expect(client.Reconnect()).NotTo(HaveOccurred())

		By("reading from the server")
		Expect(<-client.Reader()).NotTo(BeNil())

		Expect(ircd.Close()).NotTo(HaveOccurred())
	})

	It("should not be able to send empty messages", func() {
		By("creating the server")
		ircd, err := net.Listen("tcp", "127.0.0.1:5558")
		Expect(err).NotTo(HaveOccurred())

		By("starting the handler")
		result := make(chan bool)
		go IRCDaemonConn(ircd, result)

		By("connecting to the server")
		cfg.Port = 5558
		client := NewIRC(cfg)
		Expect(client.Connect()).NotTo(HaveOccurred())
		Expect(<-result).To(BeTrue())

		By("sending an empty message")
		client.Writer() <- ""

		Expect(<-client.Errors()).NotTo(BeNil())

		Expect(ircd.Close()).NotTo(HaveOccurred())
	})
})
