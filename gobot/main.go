package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/savaki/gobot"
	"github.com/savaki/gobot/builtin/listeners/slackbot"
	"github.com/savaki/gobot/builtin/providers/gocd"
)

var (
	flagSlack = cli.BoolFlag{"slack", "enable slack listener", "GOBOT_SLACK"}
)

func main() {
	app := cli.NewApp()
	app.Name = "gobot"
	app.Usage = "ThoughtWork Go plugin for chatops"
	app.Flags = []cli.Flag{
		flagSlack,
	}
	app.Action = Run
	app.Run(os.Args)
}

func assert(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func Run(c *cli.Context) {
	var echo gobot.ReceiverFunc = func(text string) (string, *gobot.Attachment, bool) {
		return "echo => " + text, nil, true
	}

	receivers := []gobot.Receiver{
		gocd.New(),
	}

	var wg sync.WaitGroup

	// start the slack listener
	if c.Bool(flagSlack.Name) {
		fmt.Println("starting slack listener")
		wg.Add(1)
		go func() {
			defer wg.Done()

			b, err := slackbot.New()
			assert(err)

			b = slackbot.WithReceiver(b, receivers...)
			err = b.Listen()
			assert(err)
		}()
	}

	wg.Wait()

}
