package gocd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/savaki/goapi"
	"github.com/savaki/gobot"
)

type receiver struct {
	api    *goapi.Client
	target *gobot.ReceiverGroup
}

func New() gobot.Handler {
	codebase := os.Getenv("GO_CODEBASE")
	username := os.Getenv("GO_USERNAME")
	password := os.Getenv("GO_PASSWORD")

	// load base configuration
	if codebase == "" {
		log.Fatalln(fmt.Errorf("ERROR - GO_CODEBASE environment variable not defined"))
	}

	g := goapi.New(codebase)
	if username != "" && password != "" {
		g = goapi.WithAuth(g, username, password)
	}

	handlers := gobot.Handlers{}
	handlers.AddFunc(handleAllBuilds)

	r := &receiver{
		api:    g,
		target: &gobot.ReceiverGroup{},
	}
	r.target.AddFunc(r.allBuilds, r.failingBuilds)

	return r
}

func (r *receiver) OnMessage(text string) (string, *gobot.Attachment, bool) {
	return r.target.OnMessage(text)
}

func handleAllBuilds(c *gobot.Context) {
}

func (r *receiver) allBuilds(text string) (string, *gobot.Attachment, bool) {
	if text != "all builds" {
		return "", nil, false
	}

	projects, err := r.api.BuildStatus()
	if err != nil {
		gobot.WrapError(err)
	}

	lines := []string{}
	for i, p := range projects {
		text := fmt.Sprintf("%d. %s => %s", i+1, p.Name, p.LastBuildStatus)
		lines = append(lines, text)
	}

	return strings.Join(lines, "\n"), nil, true
}

func (r *receiver) failingBuilds(text string) (string, *gobot.Attachment, bool) {
	if text != "failing builds" {
		return "", nil, false
	}

	projects, err := r.api.BuildStatus()
	if err != nil {
		gobot.WrapError(err)
	}

	failed := goapi.Filter.Failed.Filter(projects)
	if len(failed) == 0 {
		return "all builds green - yippee!", nil, true
	}

	lines := []string{}
	for i, p := range failed {
		text := fmt.Sprintf("%d. %s => %s", i+1, p.Name, p.LastBuildStatus)
		lines = append(lines, text)
	}

	return strings.Join(lines, "\n"), nil, true
}
