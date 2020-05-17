package verbs

import(
	"fmt"
	"time"
)

// History is intended not for users to interact with, but more as a
// backend audit log that only the framework (via VerbContext) can add
// to, so the framework has unholy knowledge of this verb (i.e. a
// layering violation). But it's nice to keep it as an actual verb, as
// we get a simple admin CLI for free.

// History is a stateful verb.
type History struct {
	Events []Event
}

type Event struct {
	Type string
	Time time.Time
	User string
	Action string
}
func (e Event)String() string {
	return fmt.Sprintf("%s [%s] %s", e.Time.Format("01/02 15:04:05 MST"), e.User, e.Action)
}
	
func (h *History)Help() string { return " " }

func NewHistory() History {
	return History{Events:[]Event{}}
}

func (h *History)Process(vc VerbContext, args []string) string {
	if h.Events == nil {
		h.Events = []Event{}
	}

	if len(args) == 0 {
		str := ""
		for _,e := range h.Events {
			str += fmt.Sprintf("%s\n", e)
		}
		return str

	} else {
		switch args[0] {
		case "-flush":
			h.Events = []Event{}
		}
	}

	return ""
}
