// Description:
//   Interact with your ThoughtWorks Go Continuous Delivery server
//
// Dependencies:
//   None
//
// Configuration:
//   GOBOT_GO_CODEBASE
//   GOBOT_GO_USERNAME
//   GOBOT_GO_PASSWORD
//
// Commands:
//   gobot go b <pipeline> - builds the pipeline specified by pipeline. List pipelines to get the list of pipelines.
//   gobot go build <pipeline> - builds the specified Go pipeline
//   gobot go list - lists Go pipelines
//   gobot go last <pipeline> - Details about the last build for the specified Go pipeline
//   gobot go status - lists failing builds

//
// Author:
//   Matt Ho

package gocd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/savaki/goapi"
	"github.com/savaki/gobot"
)

func Provider() *gobot.Provider {
	// use environment variables to instantiate the goapi
	api, err := goapi.FromEnv()
	if err != nil {
		log.Infof("Unable to load Go provider.  Go grammars will not be available. => %s", err.Error())
		return nil
	}

	// associate all our commands with the handler

	r := &receiver{api: api}
	return &gobot.Provider{
		Name: "go",
		Commands: []gobot.Command{
			{
				Grammars: []string{`go b (\S+)`, `go build (\S+)`},
				Summary:  "schedule a pipeline to run",
				Action:   r.scheduledPipeline,
			},
			{
				Grammar: "go list",
				Summary: "list all pipelines",
				Action:  r.listPipelines,
			},
			{
				Grammar: `go last (\S+)`,
				Summary: "last build status for specified pipeline",
				Action:  r.lastStatus,
			},
			{
				Grammar: "go status",
				Summary: "lists failed builds",
				Action:  r.failedBuilds,
			},
		},
	}
}

type receiver struct {
	api *goapi.Client
}

func apiFromEnv() (*goapi.Client, error) {
	// attempt to instantiate a client
	codebase := os.Getenv("GOBOT_GO_CODEBASE")
	if codebase == "" {
		return nil, fmt.Errorf("GOBOT_GO_CODEBASE environment variable not defined")
	}
	client := goapi.New(codebase)

	// associate a username and password if provided
	username := os.Getenv("GOBOT_GO_USERNAME")
	password := os.Getenv("GOBOT_GO_PASSWORD")
	if username != "" && password != "" {
		client = goapi.WithAuth(client, username, password)
	}

	return client, nil
}

func (r *receiver) listPipelines(c *gobot.Context) {
	log.WithField("provider", "gocd").Debugf("#listPipelines")

	groups, err := r.api.PipelineGroups()
	if err != nil {
		c.Fail(err)
		return
	}
	response := c.Respond("Piplines:")

	for _, g := range groups {
		for i, p := range g.Pipelines {
			response.Append(fmt.Sprintf(" %d. %s", i+1, p.Name))
		}
	}
}

func (r *receiver) scheduledPipeline(c *gobot.Context) {
	log.WithField("provider", "gocd").Debugf("#allBuilds")

	pipeline := c.Match(1)
	if err := r.api.PipelineSchedule(pipeline); err != nil {
		c.Fail(err)
		return
	}

	c.Respond(fmt.Sprintf("Scheduled pipeline, %s", pipeline))
}

func (r *receiver) lastStatus(c *gobot.Context) {
	log.WithField("provider", "gocd").Debugf("#lastStatus")

	projects, err := r.api.BuildStatus()
	if err != nil {
		c.Fail(err)
	}

	pipeline := c.Match(1)
	filtered := []goapi.Project{}
	for _, p := range projects {
		if parts := strings.Split(p.Name, " :: "); len(parts) != 2 {
			continue
		} else if parts[0] != pipeline {
			continue
		}
		filtered = append(filtered, p)
	}

	if len(filtered) == 0 {
		c.Respond(fmt.Sprintf("Unable to find a pipeline with name, %s", pipeline))
		return
	}

	response := c.Respond(fmt.Sprintf("%s:", pipeline))
	for i, p := range filtered {
		text := fmt.Sprintf("%d. %s => %s", i+1, p.Name, p.LastBuildStatus)
		response.Append(text)
	}
}

func (r *receiver) failedBuilds(c *gobot.Context) {
	log.WithField("provider", "gocd").Debugf("#failedBuilds")

	projects, err := r.api.BuildStatus()
	if err != nil {
		c.Fail(err)
	}

	failed := goapi.OnlyFailedBuilds(projects)

	if len(failed) == 0 {
		c.Respond("All builds running green!")
		return
	}

	response := c.Respond("Failed builds:")
	for i, p := range failed {
		text := fmt.Sprintf("%d. %s => %s", i+1, p.Name, p.LastBuildStatus)
		response.Append(text)
	}
}
