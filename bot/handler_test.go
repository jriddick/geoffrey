package bot

import (
	"testing"

	"github.com/jriddick/geoffrey/msg"

	. "github.com/smartystreets/goconvey/convey"
)

var testHandler = Handler{
	Name:        "Test",
	Description: "Testing handler",
	Event:       "PING",
	Run: func(bot *Bot, msg *msg.Message) (bool, error) {
		return true, nil
	},
}

func TestHandler(t *testing.T) {
	Convey("With handler", t, func() {
		Convey("Should be able to register new handler", func() {
			So(RegisterHandler(testHandler), ShouldBeNil)
		})

		Convey("Should not be able to register the same handler twice", func() {
			So(RegisterHandler(testHandler), ShouldNotBeNil)
		})
	})
}
