package irc

import "fmt"

// Color codes for IRC messages
const (
	White = iota
	Black
	Blue
	Green
	Red
	Brown
	Purple
	Orange
	Yellow
	LightGreen
	Teal
	LightCyan
	LightBlue
	Pink
	Grey
	LightGrey
	Default
)

// Foreground sets foreground color of the text to the given color
func Foreground(text string, code int) string {
	return fmt.Sprintf("\x03%d​%s\x03", code, text)
}

// Bold makes the text bold
func Bold(text string) string {
	return fmt.Sprintf("\x02%s\x02", text)
}
