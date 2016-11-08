package msg

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	maxLength = 1024
)

// Tags represents IRCv3.2 Message Tags that follows
// the psuedo-BNF below. Vendor is handled by not splitting
// so the saved key includes the vendor prefix.
//
// <tag>           ::= <key> ['=' <escaped value>]
// <key>           ::= [ <vendor> '/' ] <sequence of letters, digits, hyphens (`-`)>
// <escaped value> ::= <sequence of any characters except NUL, CR, LF, semicolon (`;`) and SPACE>
// <vendor>        ::= <host>
type Tags map[string]string

// Prefix represents an IRC Message Prefix and it follows
// the psuedo-BNF below. Nick- or servername is stored in
// the Name field.
//
// <prefix>   ::= <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
type Prefix struct {
	Name string
	User string
	Host string
}

// Message represents an IRC Message and it follow the
// psuedo-BNF below.
//
// <message>       ::= ['@' <tags> <SPACE>] [':' <prefix> <SPACE> ] <command> <params> <crlf>
// <tags>          ::= <tag> [';' <tag>]*
// <tag>           ::= <key> ['=' <escaped value>]
// <key>           ::= [ <vendor> '/' ] <sequence of letters, digits, hyphens (`-`)>
// <escaped value> ::= <sequence of any characters except NUL, CR, LF, semicolon (`;`) and SPACE>
// <vendor>        ::= <host>
// <prefix>        ::= <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
// <command>       ::= <letter> { <letter> } | <number> <number> <number>
// <SPACE>         ::= ' ' { ' ' }
// <params>        ::= <SPACE> [ ':' <trailing> | <middle> <params> ]
// <middle>        ::= <Any *non-empty* sequence of octets not including SPACE
//                     or NUL or CR or LF, the first of which may not be ':'>
// <trailing>      ::= <Any, possibly *empty*, sequence of octets not including
//                     NUL or CR or LF>
// <crlf>          ::= CR LF
type Message struct {
	Tags     Tags     `json:",omitempty"`
	Prefix   *Prefix  `json:",omitempty"`
	Command  string   `json:",omitempty"`
	Params   []string `json:",omitempty"`
	Trailing string   `json:",omitempty"`
}

func trim(r rune) bool {
	return r == '\n' || r == '\r'
}

// Bytes return the IRC message as a byte buffer
func (m *Message) Bytes() []byte {
	buf := new(bytes.Buffer)

	if m.Tags != nil {
		buf.WriteRune('@')
		length := len(m.Tags)
		curr := 0

		for key, val := range m.Tags {
			buf.WriteString(key)
			curr++

			if val != "" {
				buf.WriteRune('=')
				buf.WriteString(val)
			}

			if curr != length {
				buf.WriteRune(';')
			}
		}

		buf.WriteRune(' ')
	}

	if m.Prefix != nil {
		buf.WriteRune(':')
		buf.WriteString(m.Prefix.Name)
		if len(m.Prefix.User) > 0 {
			buf.WriteRune('!')
			buf.WriteString(m.Prefix.User)
		}
		if len(m.Prefix.Host) > 0 {
			buf.WriteRune('@')
			buf.WriteString(m.Prefix.Host)
		}
		buf.WriteRune(' ')
	}

	buf.WriteString(m.Command)

	if len(m.Params) > 0 {
		buf.WriteRune(' ')
		buf.WriteString(strings.Join(m.Params, " "))
	}

	if len(m.Trailing) > 0 {
		buf.WriteRune(' ')
		buf.WriteRune(':')
		buf.WriteString(m.Trailing)
	}

	if buf.Len() > maxLength-2 {
		buf.Truncate(maxLength - 2)
	}

	buf.WriteRune('\r')
	buf.WriteRune('\n')

	return buf.Bytes()
}

// String returns the IRC message as a string
func (m *Message) String() string {
	return strings.TrimFunc(string(m.Bytes()), trim)
}

// ParseMessage takes an IRC message and parses
// it into a Message struct.
func ParseMessage(raw string) (*Message, error) {
	// Make sure it does not exceed max length
	if len(raw) > maxLength {
		return nil, fmt.Errorf("message is too long")
	}

	// Make sure its not empty
	if len(strings.TrimSpace(raw)) == 0 {
		return nil, fmt.Errorf("message is empty")
	}

	// Create the message
	message := new(Message)

	// Check if we have found a tag
	if raw[0] == '@' {
		// Create the tag container
		message.Tags = make(map[string]string)

		// Find the end of the tag field
		tagEnd := strings.IndexRune(raw, ' ')

		if tagEnd == -1 {
			return nil, fmt.Errorf("message ends with tags")
		}

		if tagEnd == 1 {
			return nil, fmt.Errorf("empty tag field")
		}

		tags := strings.Split(raw[1:tagEnd], ";")

		for _, tag := range tags {
			equal := strings.IndexRune(tag, '=')

			if equal > -1 {
				message.Tags[tag[:equal]] = tag[equal+1:]
			} else {
				message.Tags[tag] = ""
			}
		}

		// Remove the tags from the string
		raw = raw[tagEnd+1:]
	}

	if raw[0] == ':' {
		// Create the prefix
		message.Prefix = &Prefix{}

		// Find the end
		prefixEnd := strings.IndexRune(raw, ' ')

		if prefixEnd == -1 {
			return nil, fmt.Errorf("message ends with prefix")
		}

		if prefixEnd == 1 {
			return nil, fmt.Errorf("empty prefix")
		}

		// Get the substring
		prefix := raw[1:prefixEnd]

		// Search for the user
		userToken := strings.IndexRune(prefix, '!')

		// Check if we found the user
		if userToken > -1 {
			// Set the name prefix
			message.Prefix.Name = prefix[:userToken]

			// Update the prefix
			prefix = prefix[userToken+1:]

			// See if we can find the host
			hostToken := strings.IndexRune(prefix, '@')

			if hostToken > -1 {
				message.Prefix.User = prefix[:hostToken]
				message.Prefix.Host = prefix[hostToken+1:]
			} else {
				message.Prefix.User = prefix
			}
		} else {
			// See if we can find the host
			hostToken := strings.IndexRune(prefix, '@')

			if hostToken > -1 {
				message.Prefix.Name = prefix[:hostToken]
				message.Prefix.Host = prefix[hostToken+1:]
			} else {
				message.Prefix.Name = prefix
			}
		}

		// Remove the prefix
		raw = raw[prefixEnd+1:]
	}

	// Search for the end of command
	commandEnd := strings.IndexRune(raw, ' ')

	if commandEnd > -1 {
		// Set the command
		message.Command = raw[:commandEnd]

		// Find the trailing
		trailingToken := strings.Index(raw, " :")

		// Check if we found the token
		if trailingToken > -1 {
			// Set the trailing
			message.Trailing = strings.TrimFunc(raw[trailingToken+2:], trim)

			if commandEnd+1 < trailingToken-1 {
				// Set the params
				message.Params = strings.Split(strings.TrimSpace(raw[commandEnd+1:trailingToken]), " ")
			}
		} else {
			// Set the params
			message.Params = strings.Split(raw[commandEnd+1:], " ")
		}
	} else {
		// Set the command
		message.Command = strings.TrimSpace(raw)
	}

	return message, nil
}
