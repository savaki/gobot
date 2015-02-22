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
			So(command.OnLoad(), ShouldBeNil)

			grammar, matches, ok := command.matcher.match("hello world")
			So(grammar, ShouldEqual, command.Grammar)
			So(matches, ShouldResemble, []string{"hello world", "hello", "wor"})
			So(ok, ShouldBeTrue)
		})
	})

	Convey("Given a command with pattern matching", t, func() {
		command := &Command{
			Grammar: `go b (\S+)`,
			Action:  action,
		}

		Convey("Then I expect the matcher to match", func() {
			So(command.OnLoad(), ShouldBeNil)

			grammar, matches, ok := command.matcher.match("go b FirstPipeline")
			So(grammar, ShouldEqual, command.Grammar)
			So(matches, ShouldResemble, []string{"go b FirstPipeline", "FirstPipeline"})
			So(ok, ShouldBeTrue)
		})
	})

	Convey("Given a command", t, func() {
		command := &Command{
			Grammar: "hello world",
			Action:  action,
		}

		Convey("When I call init", func() {
			err := command.OnLoad()

			Convey("Then I expect no errors", func() {
				So(err, ShouldBeNil)
			})

			Convey("And I expect things to be initialized", func() {
				So(command.matcher, ShouldNotBeNil)

				grammar, matches, ok := command.matcher.match(command.Grammar)
				So(grammar, ShouldEqual, command.Grammar)
				So(matches, ShouldResemble, []string{command.Grammar})
				So(ok, ShouldBeTrue)
			})
		})

		Convey("When I call #handle", func() {
			err := command.OnLoad()
			So(err, ShouldBeNil)

			ctx := &Context{Text: command.Grammar}
			resp, ok := command.OnMessage(ctx)

			Convey("Then I expect our command to be executed", func() {
				So(ok, ShouldBeTrue)
				So(resp.Text, ShouldEqual, content)
			})
		})
	})
}
