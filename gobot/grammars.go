package main

import (
	"fmt"
	"sort"

	"github.com/savaki/gobot"
)

func allGrammars(name string, handler gobot.Handler) gobot.Handler {
	g := "grammars"
	return &gobot.Command{
		Grammar: g,
		Action: func(c *gobot.Context) {
			response := c.Respond("Supported grammars:")

			// 1. retrieve all the examples
			all := handler.Examples()
			all = append(all, gobot.Example{
				Provider: BuiltinProvider,
				Grammar:  g,
				Summary:  "list all grammars",
			})

			// 2. determine the list of providers
			providers := all.Providers()
			sort.Strings(providers)

			// 3. render the grammars by provider
			grouped := all.GroupBy(func(e gobot.Example) string { return e.Provider })
			for _, p := range providers {
				examples := grouped[p]
				sort.Sort(examples)

				response.Append("")
				response.Append(fmt.Sprintf("[%s]", p))
				for _, e := range examples {
					response.Append(fmt.Sprintf("* %s %s - %s", name, e.Grammar, e.Summary))
				}
			}
		},
	}
}
