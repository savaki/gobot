package slackbot

import (
	"fmt"
	"os"
	"regexp"

	"github.com/savaki/gobot"
	"github.com/savaki/slack"
)

const (
	DefaultName = "gobot"
)

var (
	receivers = []gobot.Receiver{}
)

func New() (*Listener, error) {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("ERROR - missing env variable, SLACK_TOKEN")
	}

	name := os.Getenv("SLACK_NAME")
	if name == "" {
		name = DefaultName
	}

	api := slack.New(token)

	pattern := fmt.Sprintf(`\s*%s\s+(.*)$`, name)
	matcher, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &Listener{
		Matcher:   matcher,
		Receivers: []gobot.Receiver{},
		api:       api,
		Name:      name,
	}, nil
}

func WithReceiver(l *Listener, r ...gobot.Receiver) *Listener {
	l.Receivers = append(l.Receivers, r...)
	return l
}

type Listener struct {
	Matcher   *regexp.Regexp
	Receivers []gobot.Receiver
	api       *slack.Client
	Name      string
}

func (l *Listener) Listen() error {
	return l.api.Listen(l)
}

func (c *Listener) OnMessage(event slack.MessageEvent) error {
	if matches := c.Matcher.FindStringSubmatch(event.Text); len(matches) > 1 {
		text := matches[1]
		if resp, attachment, match := c.processMessage(text); match {
			c.respond(event, resp, attachment)
		}
	}

	return nil
}

func (c *Listener) respond(event slack.MessageEvent, text string, attachment *gobot.Attachment) error {
	_, err := c.api.PostMessage(slack.PostMessageReq{
		Channel:  event.Channel,
		Text:     text,
		Username: c.Name,
	})
	return err
}

func (c *Listener) processMessage(text string) (string, *gobot.Attachment, bool) {
	if c.Receivers != nil {
		for _, r := range c.Receivers {
			if resp, attachment, match := r.OnMessage(text); match {
				return resp, attachment, match
			}
		}
	}

	return "", nil, false
}
