package verbs

import(
	"fmt"
	"strings"
	"time"
)

type IWillActions struct {
	Actions []IWillAction
}

type IWillAction struct {
	Time time.Time
	User string
	Action string
}
func (a IWillAction)String() string {
	return fmt.Sprintf("%s [%s] will %s", a.Time.Format("01/02 15:04:05 MST"), a.User, a.Action)
}
	
func (i *IWillActions)Help() string { return "will eat cheese for money" }

func (i *IWillActions)Process(vc VerbContext, args []string) string {
	if len(args) < 1 { return i.Help() }
	
	if i.Actions == nil {
		i.Actions = []IWillAction{}
	}
	
	switch args[0] {
	case "-flush":
		i.Actions = []IWillAction{}
		
	case "will":
		if vc.User == "" {
			return "no you won't, I don't know you"
		} else if len(args) == 1 {
			return "so you claim, pah"
		}
		new := strings.Join(args[1:], " ")
		i.Actions = append(i.Actions, IWillAction{time.Now(), vc.User, new})
		return "thy will be done"

	case "list":
		str := ""
		for _,a := range i.Actions {
			str += fmt.Sprintf("%s\n", a)
		}
		return str
		
	default:
		return i.Help()
	}

	return ""
}
