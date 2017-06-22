package irc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestColors(t *testing.T) {
	Convey("With colors package", t, func() {
		Convey("It should be able to colorize a string", func() {
			result := Foreground("red", Red)
			So(result, ShouldEqual, "\x034red\x03")
		})

		Convey("It should be able to bold a string", func() {
			result := Bold("bold")
			So(result, ShouldEqual, "\x02bold\x02")
		})
	})
}
