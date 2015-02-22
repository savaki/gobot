package slackbot

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/savaki/gobot"
	"github.com/savaki/slack"
)

const (
	DefaultName = "gobot"
)

func Listen(name string, handler gobot.Handler) error {
	log.WithField("provider", "slackbot").Debugf("starting slack listener with name, %s", name)

	// 1. retrieve the api
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		return fmt.Errorf("ERROR - missing env variable, SLACK_TOKEN")
	}
	api := slack.New(token)

	// 2. create a matcher for the name
	pattern := fmt.Sprintf(`\s*%s\s+(.*)$`, name)
	matcher, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	r := robot{
		api:     api,
		name:    name,
		matcher: matcher,
		handler: handler,
	}

	// 3. pass to the api to listen
	return api.Listen(r)
}

type robot struct {
	api     *slack.Client
	name    string
	matcher *regexp.Regexp
	handler gobot.Handler
}

func (r robot) OnMessage(event slack.MessageEvent) error {
	log.WithField("provider", "slackbot").Debugf("[RAW] => %s", event.Text)
	if matches := r.matcher.FindStringSubmatch(event.Text); len(matches) > 1 {
		text := strings.TrimSpace(matches[1])

		log.WithField("provider", "slackbot").Debugf("[IN]  => %s", text)
		ctx := &gobot.Context{
			User: event.User,
			Text: text,
		}
		if response, ok := r.handler.OnMessage(ctx); ok {
			r.respond(event, response)
		}
	}

	return nil
}

func (r robot) respond(event slack.MessageEvent, response *gobot.Response) error {
	if log.GetLevel() == log.DebugLevel {
		text := response.Text
		if i := strings.Index(text, "\n"); i > 0 {
			text = text[0:i] + "..."
		}
		log.WithField("provider", "slackbot").Debugf("[OUT] => %s", text)
	}

	// send text messages
	if response.Text != "" {
		err := r.respondText(event.Channel, response.Text)
		if err != nil {
			return err
		}
	}

	// upload attachments
	if response.Attachments != nil {
		for _, a := range response.Attachments {
			req := &slack.FilesUploadReq{
				Content:  a.Content,
				Filetype: a.ContentType,
				Filename: a.Filename,
				Title:    a.Title,
				Channels: []string{event.Channel},
			}
			resp, err := r.api.FilesUpload(req)
			if err != nil {
				return err
			}

			if !resp.Ok {
				r.respondText(event.Channel, resp.Error)
				return errors.New(resp.Error)
			}
		}
	}

	return nil
}

func (r robot) respondText(channel, text string) error {
	_, err := r.api.PostMessage(slack.PostMessageReq{
		Channel:  channel,
		Text:     text,
		Username: r.name,
	})
	return err
}
