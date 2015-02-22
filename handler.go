package gobot

type Handler interface {
	// return grammar examples
	Examples() Examples

	// Called once when the handler is first loaded
	OnLoad() error

	// Called when each message arrives
	OnMessage(string) (*Response, bool)
}

// -------------------------------------------------------

type Handlers []Handler

func (h Handlers) Examples() Examples {
	examples := Examples{}

	for _, handler := range h {
		for _, example := range handler.Examples() {
			examples = append(examples, example)
		}
	}

	return examples
}

func (h Handlers) WithHandlers(handlers ...Handler) Handlers {
	return append(h, handlers...)
}

func (h Handlers) WithCommands(commands ...*Command) Handlers {
	handlers := make([]Handler, len(commands))
	for i, c := range commands {
		handlers[i] = c
	}
	return h.WithHandlers(handlers...)
}

func (h Handlers) WithProvider(provider *Provider) Handlers {
	if provider != nil {
		return h.WithHandlers(provider.asHandlers()...)
	}
	return h
}

func (h Handlers) OnLoad() error {
	for _, handler := range h {
		if err := handler.OnLoad(); err != nil {
			return err
		}
	}

	return nil
}

func (h Handlers) OnMessage(text string) (*Response, bool) {
	for _, handler := range h {
		if resp, ok := handler.OnMessage(text); ok {
			return resp, ok
		}
	}

	return nil, false
}
