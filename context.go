package gobot

import log "github.com/Sirupsen/logrus"

// -------------------------------------------------------

type Context struct {
	User     string
	Text     string
	matches  []string
	response *Response
	ok       bool
}

func (c *Context) Match(index int) string {
	if index > len(c.matches) {
		log.WithField("grammar", "match-err").Warnf("invalid #Match(%d) request => %v [%d]", index, c.matches[0], len(c.matches))
		return ""
	}

	return c.matches[index]
}

func (c *Context) Upload(attachment Attachment) {
	c.ok = true

	if c.response == nil {
		c.response = &Response{}
	}

	response := c.response
	if response.Attachments == nil {
		response.Attachments = []Attachment{}
	}

	response.Attachments = append(response.Attachments, attachment)
}

func (c *Context) Respond(text string) *Response {
	c.ok = true

	if c.response == nil {
		c.response = &Response{}
	}

	c.response.Text = text
	return c.response
}

func (c *Context) Fail(err error) {
	c.Respond(err.Error())
}

// -------------------------------------------------------

type Response struct {
	Text        string
	Attachments []Attachment
}

func (r *Response) Append(text string) *Response {
	r.Text = r.Text + "\n" + text
	return r
}
