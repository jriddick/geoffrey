package irc_test

import (
	. "github.com/jriddick/geoffrey/irc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Irc", func() {
	var (
		client *IRC
		cfg    Config
		emp    Config
	)

	BeforeEach(func() {
		cfg = Config{
			Hostname:           "localhost",
			Port:               6697,
			Secure:             false,
			InsecureSkipVerify: false,
		}
		emp = Config{}
		client = NewIRC(cfg)
	})

	AfterEach(func() {
		client.Disconnect()
	})

	It("should be able to create an irc client", func() {
		Expect(NewIRC(cfg)).NotTo(BeNil())
	})

	It("should be able to connect", func() {
		Expect(client.Connect()).NotTo(HaveOccurred())
	})

	It("should not be able to connect with an empty configuration", func() {
		Expect(NewIRC(emp).Connect()).To(HaveOccurred())
	})

	It("should be able to read from the server", func(done Done) {
		By("connecting to the server")
		Expect(client.Connect()).NotTo(HaveOccurred())

		By("getting the reader")
		reader := client.Reader()
		Expect(reader).NotTo(BeNil())

		Expect(<-reader).NotTo(BeNil())
		close(done)
	}, 2)

	It("should be able to register to the server", func(done Done) {
		By("connecting to the server")
		Expect(client.Connect()).NotTo(HaveOccurred())

		By("getting the reader and writer")
		reader := client.Reader()
		Expect(reader).NotTo(BeNil())
		writer := client.Writer()
		Expect(writer).NotTo(BeNil())

		By("sending nick and user")
		writer <- "NICK geoffrey"
		writer <- "USER geoffrey 0 * :geoffrey"

		for {
			line, ok := <-reader
			Expect(ok).To(BeTrue())

			if line.Command == RPL_WELCOME {
				break
			}
		}

		close(done)
	}, 2)

	It("should be able to reconnect to the server", func(done Done) {
		By("connecting to the server")
		Expect(client.Connect()).NotTo(HaveOccurred())

		By("reading from the server")
		Expect(<-client.Reader()).NotTo(BeNil())

		By("reconnecting")
		Expect(client.Reconnect()).NotTo(HaveOccurred())

		By("reading from the server")
		Expect(<-client.Reader()).NotTo(BeNil())

		close(done)
	}, 2)
})
