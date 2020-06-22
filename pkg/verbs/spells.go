package verbs

import(
	"fmt"
	"strconv"
	
	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)

// Spells is a stateless verb, since it operates on the character
// state found in the context
type Spells struct {}


// Spellbook stuff is optional; if you leave it empty, then you know all spells.
//   spells book add magic-missile  # add it to spellbook
//   spells book                    # print out spellbook


// spells set 1.2 cure-wounds     # fills up spellslot 1.2 with the cure-wounds spell
// spells cast 1.1                # cast it !
// spells resetall                # re-memorize all the spells

func (s Spells)Help() string { return "[set 1.3 cure-wounds] [cast 1.3] [resetall]" }

func (s Spells)Process(vc VerbContext, args []string) string {
	if vc.Character == nil { return "no character loaded :(" }
	if len(args) < 1 { return vc.Character.SpellSlots.String() }

	switch args[0] {
	case "-flush":
		vc.Character.Spellbook = character.NewSpellbook()
		vc.Character.SpellSlots = character.NewSpellSlots()
		return "(flushed)"

	case "-init":
		if len(args) < 3 { return "-init KIND 5 3 1" }
		max,_ := Atois(args[2:])
		vc.Character.Init(max, args[1])
		return "ok"
/*
	case "book":
		switch len(args) {
		case 1: return vc.Character.Spellbook.String()
		case 3: return "book add is TBD"
		default: return s.Help()
		}
*/

	case "set":
		if len(args) != 3 { return s.Help() }

		if err := vc.Character.SpellSlots.Memorize(&vc.Character.Spellbook, args[1], args[2]); err != nil {
			return "you dunce - " + err.Error()
		}
		return "ooooh"

	case "cast":
		if len(args) != 2 { return s.Help() }
		if spell,err := vc.Character.SpellSlots.Cast(args[1]); err != nil {
			return err.Error()
		} else if spell == "" {
			return "*fizzle*"
		} else {
			vc.LogEvent("cast " + spell)

			str := fmt.Sprintf("%s casts '%s'", vc.User, spell)
			if s := rules.TheRules.SpellList.LookupFirst(spell); s != nil {
				str += "\n\n" + s.Description()
			}
			return str
		}

	case "resetall":
		vc.Character.SpellSlots.Refresh()
		vc.LogEvent("refreshed their spell slots")
		return "aah, that's better"
		
	default:
		return s.Help()
	}

	return ""
}

// Atois does Atoi across a slice of words
func Atois(in []string) ([]int, error) {
	out := []int{}
	for _,s := range in {
		if i,err := strconv.Atoi(s); err != nil {
			return nil, err
		} else {
			out = append(out, i)
		}
	}
	return out, nil
}
