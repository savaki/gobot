package gobot

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommand(t *testing.T) {
	content := "hello world"
	action := func(c *Context) {
		c.Respond(content)
	}

	Convey("Given a command with pattern matching", t, func() {
		command := &Command{
			Grammar: "(hello) (wor)ld",
			Action:  action,
		}

		Convey("Then I expect the matcher to match", func() {
			So(command.init(), ShouldBeNil)
			So(command.matcher.FindStringSubmatch("hello world"), ShouldResemble, []string{"hello world", "hello", "wor"})
		})
	})

	Convey("Given a command", t, func() {
		command := &Command{
			Grammar: "hello world",
			Action:  action,
		}

		Convey("When I call init", func() {
			err := command.init()

			Convey("Then I expect no errors", func() {
				So(err, ShouldBeNil)
			})

			Convey("And I expect things to be initialized", func() {
				So(command.matcher, ShouldNotBeNil)

				parts := command.matcher.FindStringSubmatch(command.Grammar)
				So(parts, ShouldResemble, []string{command.Grammar})
			})
		})

		Convey("When I call #handle", func() {
			err := command.init()
			So(err, ShouldBeNil)

			resp, ok := command.OnMessage(command.Grammar)

			Convey("Then I expect our command to be executed", func() {
				So(ok, ShouldBeTrue)
				So(resp.Text, ShouldEqual, content)
			})
		})
	})
}
