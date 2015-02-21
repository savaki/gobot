package gobot

import "io"

type Attachment struct {
	Name        string
	Reader      io.ReadCloser
	ContentType string
}

type Receiver interface {
	OnMessage(text string) (string, *Attachment, bool)
}

type ReceiverFunc func(string) (string, *Attachment, bool)

func (r ReceiverFunc) OnMessage(text string) (string, *Attachment, bool) {
	return r(text)
}

type ReceiverGroup struct {
	Receivers []Receiver
}

func (rg *ReceiverGroup) Add(r ...Receiver) *ReceiverGroup {
	if rg.Receivers == nil {
		rg.Receivers = []Receiver{}
	}
	rg.Receivers = append(rg.Receivers, r...)
	return rg
}

func (rg *ReceiverGroup) AddFunc(r ...ReceiverFunc) *ReceiverGroup {
	for _, fn := range r {
		var receiver Receiver = fn
		rg.Add(receiver)
	}
	return rg
}

func (rg *ReceiverGroup) OnMessage(text string) (string, *Attachment, bool) {
	if rg.Receivers != nil {
		for _, r := range rg.Receivers {
			if resp, attachment, match := r.OnMessage(text); match {
				return resp, attachment, match
			}
		}
	}

	return "", nil, false
}

func WrapError(err error) (string, *Attachment, bool) {
	return err.Error(), nil, true
}
