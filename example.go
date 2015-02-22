package gobot

// -------------------------------------------------------

type Example struct {
	Grammar  string
	Summary  string
	Provider string
}

type Examples []Example

func (e Examples) Len() int {
	return len(e)
}

func (e Examples) Less(i, j int) bool {
	return e[i].Grammar < e[j].Grammar
}

func (e Examples) Swap(i, j int) {
	v := e[i]
	e[i] = e[j]
	e[j] = v
}

func (e Examples) Filter(f func(e Example) bool) Examples {
	examples := Examples{}
	for _, example := range e {
		if f(example) {
			examples = append(examples, example)
		}
	}
	return examples
}

func (e Examples) Providers() []string {
	providers := []string{}
	for key, _ := range e.GroupBy(func(e Example) string { return e.Provider }) {
		providers = append(providers, key)
	}
	return providers
}

func (e Examples) GroupBy(f func(e Example) string) map[string]Examples {
	groups := map[string]Examples{}

	for _, example := range e {
		key := f(example)
		if examples, found := groups[key]; found {
			groups[key] = append(examples, example)
		} else {
			groups[key] = Examples{example}
		}
	}

	return groups
}
