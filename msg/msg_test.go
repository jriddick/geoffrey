package msg_test

import (
	. "github.com/jriddick/geoffrey/msg"

	"log"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var MessageTests = [...]*struct {
	Result     *Message
	Message    string
	ShouldFail bool
}{
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "syrk",
				User: "kalt",
				Host: "millennium.stealth.net",
			},
			Command:  "QUIT",
			Trailing: "Gone to have lunch",
		},
		Message: ":syrk!kalt@millennium.stealth.net QUIT :Gone to have lunch",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "Trillian",
			},
			Command:  "SQUIT",
			Params:   []string{"cm22.eng.umd.edu"},
			Trailing: "Server out of control",
		},
		Message: ":Trillian SQUIT cm22.eng.umd.edu :Server out of control",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "WiZ",
				User: "jto",
				Host: "tolsun.oulu.fi",
			},
			Command: "JOIN",
			Params:  []string{"#Twilight_zone"},
		},
		Message: ":WiZ!jto@tolsun.oulu.fi JOIN #Twilight_zone",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "WiZ",
				User: "jto",
				Host: "tolsun.oulu.fi",
			},
			Command:  "PART",
			Params:   []string{"#playzone"},
			Trailing: "I lost",
		},
		Message: ":WiZ!jto@tolsun.oulu.fi PART #playzone :I lost",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "WiZ",
				User: "jto",
				Host: "tolsun.oulu.fi",
			},
			Command: "MODE",
			Params:  []string{"#eu-opers", "-l"},
		},
		Message: ":WiZ!jto@tolsun.oulu.fi MODE #eu-opers -l",
	},
	{
		Result: &Message{
			Command: "MODE",
			Params:  []string{"&oulu", "+b", "*!*@*.edu", "+e", "*!*@*.bu.edu"},
		},
		Message: "MODE &oulu +b *!*@*.edu +e *!*@*.bu.edu",
	},
	{
		Result: &Message{
			Command:  "PRIVMSG",
			Params:   []string{"#channel"},
			Trailing: "Message with :colons!",
		},
		Message: "PRIVMSG #channel :Message with :colons!",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "irc.vives.lan",
			},
			Command:  "251",
			Params:   []string{"test"},
			Trailing: "There are 2 users and 0 services on 1 servers",
		},
		Message: ":irc.vives.lan 251 test :There are 2 users and 0 services on 1 servers",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "irc.vives.lan",
			},
			Command:  "376",
			Params:   []string{"test"},
			Trailing: "End of MOTD command",
		},
		Message: ":irc.vives.lan 376 test :End of MOTD command",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "irc.vives.lan",
			},
			Command:  "250",
			Params:   []string{"test"},
			Trailing: "Highest connection count: 1 (1 connections received)",
		},
		Message: ":irc.vives.lan 250 test :Highest connection count: 1 (1 connections received)",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "sorcix",
				User: "~sorcix",
				Host: "sorcix.users.quakenet.org",
			},
			Command:  "PRIVMSG",
			Params:   []string{"#viveslan"},
			Trailing: "\001ACTION is testing CTCP Messages!\001",
		},
		Message: ":sorcix!~sorcix@sorcix.users.quakenet.org PRIVMSG #viveslan :\001ACTION is testing CTCP Messages!\001",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "sorcix",
				User: "~sorcix",
				Host: "sorcix.users.quakenet.org",
			},
			Command:  "NOTICE",
			Params:   []string{"midnightfox"},
			Trailing: "\001PONG 1234567890\001",
		},
		Message: ":sorcix!~sorcix@sorcix.users.quakenet.org NOTICE midnightfox :\001PONG 1234567890\001",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "a",
				User: "b",
				Host: "c",
			},
			Command: "QUIT",
		},
		Message: ":a!b@c QUIT",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "a",
				User: "b",
			},
			Command:  "PRIVMSG",
			Trailing: "Message",
		},
		Message: ":a!b PRIVMSG :Message",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "a",
				Host: "c",
			},
			Command:  "NOTICE",
			Trailing: ":::Hey!",
		},
		Message: ":a@c NOTICE ::::Hey!",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "nick",
			},
			Command:  "PRIVMSG",
			Params:   []string{"$@"},
			Trailing: "This Message contains a\ttab!",
		},
		Message: ":nick PRIVMSG $@ :This Message contains a\ttab!",
	},
	{
		Result: &Message{
			Command:  "TEST",
			Params:   []string{"$@", "", "param"},
			Trailing: "Trailing",
		},
		Message: "TEST $@  param :Trailing",
	},
	{
		Message:    ": PRIVMSG test :Invalid Message with empty prefix.",
		ShouldFail: true,
	},
	{
		Message:    ":  PRIVMSG test :Invalid Message with space prefix",
		ShouldFail: true,
	},
	{
		Result: &Message{
			Command:  "TOPIC",
			Params:   []string{"#foo"},
			Trailing: "",
		},
		Message: "TOPIC #foo",
	},
	{
		Result: &Message{
			Command:  "TOPIC",
			Params:   []string{"#foo"},
			Trailing: "",
		},
		Message: "TOPIC #foo :",
	},
	{
		Result: &Message{
			Prefix: &Prefix{
				Name: "name",
				User: "user",
				Host: "example.org",
			},
			Command:  "PRIVMSG",
			Params:   []string{"#test"},
			Trailing: "Message with spaces at the end!  ",
		},
		Message: ":name!user@example.org PRIVMSG #test :Message with spaces at the end!  ",
	},
	{
		Result: &Message{
			Command: "PASS",
			Params:  []string{"oauth:token_goes_here"},
		},
		Message: "PASS oauth:token_goes_here",
	},
	{
		Message:    "@tag=val",
		ShouldFail: true,
	},
	{
		Message:    ":loh!oh@loh.fi",
		ShouldFail: true,
	},
	{
		Message: "@tag=val;test :psyc://psyced.org/~grawity!grawity@psyced.org PRIVMSG #channel :Hello!",
		Result: &Message{
			Tags: map[string]string{
				"tag":  "val",
				"test": "",
			},
			Prefix: &Prefix{
				Name: "psyc://psyced.org/~grawity",
				User: "grawity",
				Host: "psyced.org",
			},
			Command:  "PRIVMSG",
			Params:   []string{"#channel"},
			Trailing: "Hello!",
		},
	},
	{
		Message: "@vendor/tag=val;test;one=more :nick@kcin!user!resu@host!tsoh@host PRIVMSG #channel :Hello!\r\n",
		Result: &Message{
			Tags: map[string]string{
				"vendor/tag": "val",
				"test":       "",
				"one":        "more",
			},
			Prefix: &Prefix{
				Name: "nick@kcin",
				User: "user!resu",
				Host: "host!tsoh@host",
			},
			Command:  "PRIVMSG",
			Params:   []string{"#channel"},
			Trailing: "Hello!",
		},
	},
	{
		Message:    "@ :ohloh!loh@fi.org",
		ShouldFail: true,
	},
}

type Test struct {
	Message string
	Result  *Message
}

var tests []Test

var _ = Describe("Msg", func() {
	It("ParseMessage", func() {
		for _, test := range MessageTests {
			log.Printf("%s\n", strings.TrimSpace(test.Message))
			msg, err := ParseMessage(test.Message)

			if test.ShouldFail {
				Expect(err).To(HaveOccurred())
				Expect(msg).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(msg).NotTo(BeNil())

				Expect(msg.Tags).To(Equal(test.Result.Tags))
				Expect(msg.Prefix).To(Equal(test.Result.Prefix))
				Expect(msg.Command).To(Equal(test.Result.Command))
				Expect(msg.Params).To(Equal(test.Result.Params))
				Expect(msg.Trailing).To(Equal(test.Result.Trailing))
			}
		}
	})

	It("should be able to convert IRC message to bytes", func() {
		m := &Message{
			Tags: map[string]string{
				"hello": "world",
				"money": "",
			},
			Prefix: &Prefix{
				Name: "oh",
				User: "fi!loh",
				Host: "mo@ho.org",
			},
			Command:  "PRIVMSG",
			Params:   []string{"#channel"},
			Trailing: "Hello!",
		}

		// Turn it into a string
		str := string(m.Bytes())

		// Parse the Message
		res, err := ParseMessage(str)

		// Make sure it parsed
		Expect(err).ToNot(HaveOccurred())
		Expect(res).ToNot(BeNil())

		// Make sure we get the result we expect
		Expect(res.Tags).To(Equal(m.Tags))
		Expect(res.Prefix).To(Equal(m.Prefix))
		Expect(res.Command).To(Equal(m.Command))
		Expect(res.Params).To(Equal(m.Params))
		Expect(res.Trailing).To(Equal(m.Trailing))
	})
})

func BenchmarkParseMessage_short(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ParseMessage("COMMAND arg1 :Message\r\n")
	}
}

func BenchmarkParseMessage_medium(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ParseMessage(":Namename COMMAND arg6 arg7 :Message Message Message\r\n")
	}
}

func BenchmarkParseMessage_long(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ParseMessage(":Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message Message Message Message Message\r\n")
	}
}

func BenchmarkParseMessage_max(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ParseMessage("@tag=val;tag1;tag2;tag3;tag4 :Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message Message Message Message Message\r\n")
	}
}

func BenchmarkBuildMessage_max(b *testing.B) {
	b.ReportAllocs()

	// Create a message
	m := &Message{
		Tags: map[string]string{
			"hello": "world",
			"money": "",
		},
		Prefix: &Prefix{
			Name: "oh",
			User: "fi!loh",
			Host: "mo@ho.org",
		},
		Command:  "PRIVMSG",
		Params:   []string{"#channel"},
		Trailing: "Hello!",
	}

	// Benchmark the "encoding" process
	for i := 0; i < b.N; i++ {
		m.Bytes()
	}
}
