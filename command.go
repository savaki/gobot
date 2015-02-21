package gobot

import "regexp"

type Handler interface {
	OnMessage(string) (*Response, bool)
}

type HandlerFunc func(string) (*Response, bool)

func (h HandlerFunc) OnMessage(text string) (*Response, bool) {
	return h(text)
}

// -------------------------------------------------------

type Handlers []Handler

func (h Handlers) Add(handlers ...Handler) Handler {
	h = append(h, handlers...)
	return h
}

func (h Handlers) AddFunc(handlerFuncs ...HandlerFunc) Handler {
	for _, hfn := range handlerFuncs {
		var handler Handler = hfn
		h.Add(handler)
	}
	return h
}

func (h Handlers) OnMessage(text string) (*Response, bool) {
	if h != nil {
		for _, handler := range h {
			if resp, ok := handler.OnMessage(text); ok {
				return resp, ok
			}
		}
	}

	return nil, false
}

// -------------------------------------------------------

type Command struct {
	Grammar string         `json:"grammar"`
	Run     string         `json:"run"`
	Action  func(*Context) `json:"-"`
	matcher *regexp.Regexp
}

func (c *Command) init() error {
	matcher, err := regexp.Compile(c.Grammar)
	if err != nil {
		return err
	}
	c.matcher = matcher
	return nil
}

func (c *Command) OnMessage(text string) (*Response, bool) {
	if matches := c.matcher.FindStringSubmatch(text); matches != nil {
		ctx := &Context{
			Matches: matches,
		}
		c.Action(ctx)
		return ctx.response, ctx.ok
	}

	return nil, false
}

// -------------------------------------------------------

type Context struct {
	Matches  []string
	response *Response
	ok       bool
}

func (c *Context) Respond(text string) *Response {
	c.ok = true
	c.response = &Response{Text: text}
	return c.response
}

func (c *Context) Fail() {
	c.ok = false
}

// -------------------------------------------------------

type Response struct {
	Text string
}
