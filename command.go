package gobot

import (
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// -------------------------------------------------------

type Command struct {
	Provider string         `json:"provider,omitempty"`
	Grammar  string         `json:"grammar,omitempty"`
	Grammars []string       `json:"grammars,omitempty"`
	Summary  string         `json:"summary"`
	Run      string         `json:"run"`
	Action   func(*Context) `json:"-"`
	matcher  matchers
}

func (c *Command) allGrammars() []string {
	// created a unified list of grammars
	grammars := c.Grammars
	if grammars == nil {
		grammars = []string{}
	}
	if c.Grammar != "" {
		grammars = append(grammars, c.Grammar)
	}

	return grammars
}

func (c *Command) Examples() Examples {
	examples := Examples{}

	for _, grammar := range c.allGrammars() {
		examples = append(examples, Example{
			Provider: c.Provider,
			Grammar:  grammar,
			Summary:  c.Summary,
		})
	}

	return examples
}

func (c *Command) OnLoad() error {
	// convert those into a matchers
	m := matchers{}
	for _, grammar := range c.allGrammars() {
		grammar = strings.TrimSpace(grammar)
		if grammar == "" {
			continue
		}

		grammarToCompile := grammar
		if !strings.HasPrefix(grammarToCompile, "^") {
			grammarToCompile = "^" + grammarToCompile
		}
		if !strings.HasSuffix(grammarToCompile, "$") {
			grammarToCompile = grammarToCompile + "$"
		}

		matcher, err := regexp.Compile(grammarToCompile)
		if err != nil {
			return err
		}
		m = append(m, matcherNode{
			grammar: grammar,
			matcher: matcher,
		})

		log.WithField("stage", "OnLoad").Debugf("loading grammar => %s", grammar)
	}

	c.matcher = m
	return nil
}

func (c *Command) OnMessage(ctx *Context) (*Response, bool) {
	if grammar, matches, ok := c.matcher.match(ctx.Text); ok {
		log.WithField("stage", "grammar").Debugf("'%s' matched '%s' [%d]", ctx.Text, grammar, len(matches))
		ctx.matches = matches
		c.Action(ctx)
		return ctx.response, ctx.ok
	}

	return nil, false
}

// -------------------------------------------------------

type matcherNode struct {
	grammar string
	matcher *regexp.Regexp
}

type matchers []matcherNode

func (m matchers) match(text string) (string, []string, bool) {
	for _, node := range m {
		if matches := node.matcher.FindStringSubmatch(text); matches != nil {
			return node.grammar, matches, true
		}
	}

	return "", nil, false
}

// -------------------------------------------------------

type Provider struct {
	Name     string
	Commands []Command
}

func (p *Provider) asHandlers() Handlers {
	handlers := Handlers{}

	if p.Commands != nil {
		for _, c := range p.Commands {
			command := c
			command.Provider = p.Name
			handlers = handlers.WithCommands(&command)
		}
	}

	return handlers
}
