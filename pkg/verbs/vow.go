package verbs

import(
	"fmt"
	"strings"
	"time"
)

// Vows is a stateful verb, that keep track of vows made by party members.
type Vows struct {
	Vows []Vow
}

type Vow struct {
	Time time.Time
	User string
	Action string
}
func (v Vow)String() string {
	return fmt.Sprintf("%s [%s] vowed to %s", v.Time.Format("01/02 15:04:05 MST"), v.User, v.Action)
}
	
func (v *Vows)Help() string { return "to give alms to the poor" }

func (v *Vows)Process(vc VerbContext, args []string) string {
	if len(args) < 1 { return v.Help() }
	
	if v.Vows == nil {
		v.Vows = []Vow{}
	}
	
	switch args[0] {
	case "-flush":
		v.Vows = []Vow{}
		
	case "to":
		if vc.User == "" {
			return "no you won't, I don't know you"
		} else if len(args) == 1 {
			return "so you claim, pah"
		}
		new := strings.Join(args[1:], " ")
		v.Vows = append(v.Vows, Vow{time.Now(), vc.User, new})
		return "thy will be done"

	case "list":
		str := ""
		for _,a := range v.Vows {
			str += fmt.Sprintf("%s\n", a)
		}
		return str
		
	default:
		return v.Help()
	}

	return ""
}
