package verbs

// This is an example of a stateless verb.

// HandleVerb("stateless-example", Stateless{})

import "time"

type Stateless struct{}

func (s Stateless)Help() string { return "" }

func (s Stateless)Process(vc VerbContext, args []string) string {
	return "Hi, the time is " + time.Now().String()
}
