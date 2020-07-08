package verbs

import(
	"strings"
)

// Vows is a stateless verb, in that it records the vow in the VerbContext's audit log.
type Vows struct {}
	
func (v *Vows)Help() string { return "GIVE ALMS TO THE POOR" }

func (v *Vows)Process(vc VerbContext, args []string) string {
	if len(args) < 1 { return v.Help() }
	
	if vc.User == "" {
		return "no you won't, I don't know you"
	} else if len(args) == 0 {
		return "so you claim, pah"
	}

	vc.LogEvent("vowed " + strings.Join(args, " "))
	return "thy will be done"
}
