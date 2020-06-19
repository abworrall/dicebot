package verbs

import(
	"fmt"
	"strings"
	"github.com/abworrall/dicebot/pkg/dnd5e"
)

// Rules is stateless, as the API rule objects are all magically loaded into a global var. For now.
type Rules struct{}

func (r Rules)Help() string { return "[spell like burn] [spell burning-hands]" }

func (r Rules)Process(vc VerbContext, args []string) string {
	if len(args) == 0 { return r.Help() }

	switch args[0] {
	case "spell":
		if len(args) < 2 { return r.Help() }
		if strings.EqualFold(args[1],"like") {
			if len(args) < 3 { return r.Help() }
			s := strings.Join(args[2:], " ")
			return r.ShowSpellLike(s)
		} else {
			s := strings.Join(args[1:], " ")
			return r.ShowSpell(s)
		}

	case "monster":
		return "monsters tbd"
		
	default:
		return "what...."
	}
}

func (r Rules)ShowSpellLike(s string) string {
	str := "Possible matches :-\n"
	for _,v := range dnd5e.Dnd.SpellList.Find(s) {
		str += fmt.Sprintf("[L%d (%s), %s] %s\n", v.Level, v.Class(), v.Index, v.Name)
	}
	return str
}

func (r Rules)ShowSpell(s string) string {
	if sp,exists := dnd5e.Dnd.SpellList[s]; exists {
		return sp.String()
	} else {
		return fmt.Sprintf("oooh %s sounds good, you should invent it!", s)
	}
}
