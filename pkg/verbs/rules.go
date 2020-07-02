package verbs

import(
	"fmt"
	"strings"
	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)

// Rules is a stateless verb, as the API rule objects are all
// magically loaded into a global var. For now.
type Rules struct{}

func (r Rules)Help() string { return "[spell WEB] [equip NET] [monster BAT] [buff RAGE]" }

func (r Rules)Process(vc VerbContext, args []string) string {
	if len(args) < 2 { return r.Help() }

	term := strings.Join(args[1:], " ")
	
	switch args[0] {
	case "spell":   return r.Lookup(term, rules.TheRules.SpellList)
	case "equip":   return r.Lookup(term, rules.TheRules.EquipmentList)
	case "monster": return r.Lookup(term, rules.TheRules.MonsterList)
	case "buff":    return r.Lookup(term, rules.TheRules.BuffList)

/*
	case "list":
		switch args[1] {
		case "weapons": return "TODO: implement list-by-type via Lookuper"
		}
*/

	default: return "what...."
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
