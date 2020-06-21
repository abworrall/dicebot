package verbs

import(
	"fmt"
	"strings"
	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)

// Rules is stateless, as the API rule objects are all magically loaded into a global var. For now.
type Rules struct{}

func (r Rules)Help() string { return "[spell burning-hands]" }

func (r Rules)Process(vc VerbContext, args []string) string {
	if len(args) < 2 { return r.Help() }

	term := strings.Join(args[1:], " ")
	
	switch args[0] {
	case "spell": return r.Lookup(term, rules.TheRules.SpellList)
	case "equip": return r.Lookup(term, rules.TheRules.EquipmentList)
		
	default:
		return "what...."
	}
}

func (r Rules)Lookup(s string, list rules.Lookuper) string {
	matches := list.Lookup(s)
	switch len(matches) {
	case 0: return fmt.Sprintf("couldn't find anything like '%s'", s)
	case 1: return matches[0].Description()
	default:
		str := "Possible matches :-\n"
		for _,match := range matches {
			str += match.Summary() + "\n"
		}
		return str
	}
}
