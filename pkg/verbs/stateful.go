package verbs

// A stateful verb should be autopopulated on each invocation, and persisted at the end
// (if it returns something non-empty).

type Stateful struct{
	State string
}

func NewStateful() *Stateful {
	return &Stateful{}
}

func (s *Stateful)Help() string { return "no help for you" }

func (s *Stateful)Process(vc VerbContext, args []string) string {
	if len(args) == 0 { return "[" + s.State + "]" }

	switch args[0] {
	case "set": s.State = args[1]
	case "append": s.State += args[1]
	case "reset": s.State = ""
	}
	
	return "{" + s.State + "}"
}
