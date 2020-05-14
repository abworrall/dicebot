package verbs

// This is an example of a stateful verb. The primary object
// (`Stateful`) is populated before `Process` is called, and is then
// persisted afterwards; so `Process` can always assume that state
// exists. Note, though, it will need to lazily initialize any
// substructure (e.g. maps)

// HandleVerb("stateful-example", &Stateful{})

type Stateful struct{
	State string
}

func (s *Stateful)Help() string { return "" }

func (s *Stateful)Process(vc VerbContext, args []string) string {
	if len(args) == 0 { return "[" + s.State + "]" }

	switch args[0] {
	case "set":    s.State = args[1]
	case "append": s.State += args[1]
	case "reset":  s.State = ""
	}
	
	return "{" + s.State + "}"
}
