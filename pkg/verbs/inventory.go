package verbs

import(
	"fmt"
	"strings"
)

// Inventory is a stateless verb, since it operates on the character state
// found in the context
type Inventory struct {}
	
func (i Inventory)Help() string { return "[stash item] [list] [remove N] [use N]" }

func (i Inventory)Process(vc VerbContext, args []string) string {
	if len(args) < 1 { return i.Help() }
	
	inv := vc.Character.Inventory
	inv.MaybeInit() // FIXME: this maybeinit is getting way out of hand
	
	switch args[0] {
	case "-flush":
		inv.Clear()
		
	case "stash":
		if len(args) == 1 { return "what do you want to stash, eh ?" }
		inv.Append(strings.Join(args[1:], " "))
		return "item stashed"

	case "list":
		return inv.String()

	case "use":
		if n,str := inv.ParseIndex(args[1:]); str != "" {
			return str
		} else {
			return fmt.Sprintf("%s uses their %s\n", vc.User, inv.Items[n])
		}

	case "remove":
		if n,str := inv.ParseIndex(args[1:]); str != "" {
			return str
		} else {
			inv.Remove(n)
		}

	default:
		return i.Help()
	}

	return ""
}
