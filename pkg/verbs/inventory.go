package verbs

import(
	"fmt"
	"strings"
)

// Inventory is a stateless verb, since it operates on the character state
// found in the context
type Inventory struct {}
	
func (i Inventory)Help() string { return "[stash FOO BAR] [remove NN] [use NN]" }

func (i Inventory)Process(vc VerbContext, args []string) string {
	c := vc.Character

	if len(args) < 1 { return c.Inventory.String() }

	switch args[0] {
	case "-flush":
		c.Inventory.Clear()

	case "stash":
		if len(args) == 1 { return "what do you want to stash, eh ?" }
		c.Inventory.Append(strings.Join(args[1:], " "))
		return "item stashed"

	case "remove":
		if n,str := c.Inventory.ParseIndex(args[1:]); str != "" {
			return str
		} else {
			c.Inventory.Remove(n)
		}

	case "use":
		if n,str := c.Inventory.ParseIndex(args[1:]); str != "" {
			return str
		} else {
			vc.LogEvent(fmt.Sprintf("uses their %s", c.Inventory.Items[n]))
			return fmt.Sprintf("%s uses their %s\n", vc.User, c.Inventory.Items[n])
		}

	default:
		return i.Help()
	}

	return ""
}
