package verbs

import(
	"fmt"
	"strings"
	
	"github.com/abworrall/dicebot/pkg/character"
)

// Spells is a stateless verb, since it operates on the character
// state found in the context
type SimpleSpells struct {}

// spells memorize 1:1 cure light wounds
// spells cast 1:1                         # cast it !
// spells cast 1:2                         # cast the other one !
// spells rememorize                       # re-memorize them all

func (s SimpleSpells)Help() string { return "[memorize 1:3 some magic spell] [cast 1:3] [rememorize]" }

func (s SimpleSpells)Process(vc VerbContext, args []string) string {
	if len(args) < 1 { return vc.Character.SpellSlots.String() }

	switch args[0] {
	case "-flush":
		vc.Character.SpellSlots = character.NewSpellSlots()
		return "(flushed)"

	case "-setmax":
		if len(args) < 3 { return "-setmax 5 3 1" }
		max,_ := Atois(args[1:])
		vc.Character.SetMax(max)
		return "ok"

	case "memorize":
		if len(args) < 3 { return s.Help() }
		spell := character.Spell{Name: strings.Join(args[2:], " ")}
		slotIdx := character.SlotIndex(args[1])

		if err := vc.Character.SpellSlots.MemorizeSpell(&spell, slotIdx); err != nil {
			return "you dunce - " + err.Error()
		}
		return "ooooh"

	case "cast":
		if len(args) != 2 { return s.Help() }
		slotIdx := character.SlotIndex(args[1])
		if spell,err := vc.Character.SpellSlots.Cast(slotIdx); err != nil {
			return err.Error()
		} else if spell == nil {
			return "*fizzle*"
		} else {
			vc.LogEvent("cast " + spell.String())
			return fmt.Sprintf("%s casts '%s'", vc.User, spell)
		}

	case "rememorize":
		vc.Character.SpellSlots.Refresh()
		vc.LogEvent("refreshed their spell slots")
		return "aah, that's better"
		
	default:
		return s.Help()
	}

	return ""
}
