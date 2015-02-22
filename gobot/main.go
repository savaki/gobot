package main

import (
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/savaki/gobot"
	"github.com/savaki/gobot/builtin/listeners/slackbot"
	"github.com/savaki/gobot/builtin/providers/gocd"
	"github.com/savaki/gobot/builtin/providers/mfa"
)

const (
	BuiltinProvider = "builtin"
)

var (
	flagSlack   = cli.BoolFlag{"slack", "enable slack listener", "GOBOT_SLACK"}
	flagName    = cli.StringFlag{"name", "gobot", "the name of the bot", "GOBOT_NAME"}
	flagVerbose = cli.BoolFlag{"verbose", "verbose level logging", "GOBOT_VERBOSE"}
)

func main() {
	app := cli.NewApp()
	app.Name = "gobot"
	app.Usage = "ThoughtWork Go plugin for chatops"
	app.Flags = []cli.Flag{
		flagSlack,
		flagName,
		flagVerbose,
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
	name := c.String(flagName.Name)
	if c.Bool(flagVerbose.Name) {
		log.SetLevel(log.DebugLevel)
		log.Debugf("setting log level to debug")
	}

	handlers := gobot.Handlers{}
	handlers = handlers.
		WithProvider(gocd.Provider()).
		WithProvider(mfa.Provider())
	handlers = handlers.WithHandlers(allGrammars(name, handlers))

	err := handlers.OnLoad()
	assert(err)

	var wg sync.WaitGroup

	// start the slack listener
	if c.Bool(flagSlack.Name) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := slackbot.Listen(name, handlers)
			assert(err)
		}()
	}

	wg.Wait()

}
