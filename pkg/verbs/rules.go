package verbs

import(
	"fmt"
	"strings"
	"github.com/abworrall/dicebot/pkg/dnd5e"
)

// Rules is stateless, as the API rule objects are all magicaly loaded into a global var. For now.
type Rules struct{}

func (r Rules)Help() string { return "spell [like burn] [burning-hands]" }

func (r Rules)Process(vc VerbContext, args []string) string {

	switch args[0] {
	case "spell":
		if strings.EqualFold(args[1],"like") {
			s := strings.Join(args[2:], " ")
			return r.ShowSpellLike(s)
		} else {
			s := strings.Join(args[1:], " ")
			return r.ShowSpell(s)
		}

	case "monster":
		return "monsters tbd"
		//s := strings.Join(args[1:], " ")
		//return ShowMonster(p, s)

	default:
		return "what...."
	}
}

func (r Rules)ShowSpellLike(s string) string {
	str := "Possible matches :-\n"
	for _,spell := range dnd5e.Dnd.SpellList.Find(s) {
		str += fmt.Sprintf("[L%d, %s] %s\n", spell.Level, spell.Index, spell.Name)
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
