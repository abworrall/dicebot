package verbs

import(
	"fmt"
	"strconv"
	"strings"
	
	"github.com/abworrall/dicebot/pkg/character"
)

// Spells is a stateless verb, since it operates on the character
// state found in the context
type Spells struct {}


// spells learn 1 some rubbish    # level 1, 1st spell
// spells learn 1 magic missile   # level 2, 2nd spell, etc
// spells book                    # print out the spellbook

// spells memorize 1:2 1          # memorizes the spellbook entry 1:2 into (1st level) slot 1
// spells memorize 1:2 2          # memorizes the spellbook entry 1:2 into (1st level) slot 2
// spells cast 1:1                # cast it !
// spells cast 1:2                # cast the other one !
// spells refresh                 # re-memorize

func (s Spells)Help() string { return "[learn LVL Some spell] [showbook] [memorize 1:3 2] [cast 1:3] [rememorize]" }

func (s Spells)Process(vc VerbContext, args []string) string {
	if vc.Character == nil { return "no character loaded :(" }
	if len(args) < 1 { return vc.Character.SpellSlots.String() }

	switch args[0] {
	case "-flush":
		vc.Character.Spellbook = character.NewSpellbook()
		vc.Character.SpellSlots = character.NewSpellSlots()
		return "(flushed)"

	case "-setmax":
		if len(args) < 3 { return "-setmax 5 3 1" }
		max,_ := Atois(args[1:])
		vc.Character.SetMax(max)
		return "ok"
		
	case "showbook":
		return vc.Character.Spellbook.String()

	case "learn": // args[1]. args[2..]
		if len(args) < 3 { return s.Help() }
		if lvl,err := strconv.Atoi(args[1]); err != nil || lvl < 1 || lvl > 9 {
			return "nonsense"
		} else {
			vc.Character.Spellbook.Learn(lvl, strings.Join(args[2:], " "))
			return "such learnings"
		}

	case "memorize":
		if len(args) != 3 { return s.Help() }
		spellIdx := character.SpellIndex(args[1])
		idx,_ := strconv.Atoi(args[2])
		if err := vc.Character.SpellSlots.Memorize(&vc.Character.Spellbook, spellIdx, idx); err != nil {
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
