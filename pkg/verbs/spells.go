package verbs

import(
	"fmt"
	"strconv"
	
	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/spells"
)

// Spells is a stateless verb, since it operates on the character
// state found in the context
type Spells struct {}

// spells add cure-wounds        # add the spell to your working set
// spells remove inflict-wounds  # remove the spell from your working set
// spells cast cure-wounds [NN]  # cast the spell - optionally at a higher level
// spells refresh                # reser your spell slots

func (s Spells)Help() string { return "[add WEB], [remove WEB], [cast WEB [NN]], [refresh]" }

func (s Spells)Process(vc VerbContext, args []string) string {
	if vc.Character == nil { return "no character loaded :(" }
	if len(args) < 1 {
		return vc.Character.MagicString()
	}

	switch args[0] {
	case "-flush":
		vc.Character.SpellsMemorized = spells.NewSet()
		vc.Character.Slots = spells.NewSlots()
		return "(flushed)"

	case "-init":
		if len(args) < 4 { return "-init KIND ATTR 4 2 ..." }
		kind := character.ParseAttr(args[2])
		mod := vc.Character.GetModifier(kind)
		setMax := vc.Character.Level + mod
		vc.Character.SpellsMemorized = spells.NewSet()
		vc.Character.SpellsMemorized.Max = setMax
		vc.Character.SpellsMemorized.Kind = vc.Character.Class

		slotMaxes,_ := Atois(args[3:])
		vc.Character.Slots = spells.NewSlots()
		for i,max := range slotMaxes {
			vc.Character.Slots.Max[i+1] = max
		}
		vc.Character.Slots.Reset()
		return "ok!?\n" + vc.Character.MagicString()

	case "add":
		if len(args) != 2 { return s.Help() }
		if err := vc.Character.SpellsMemorized.Add(args[1]); err != nil {
			return fmt.Sprintf("add failed: %v", err)
		}
		return "oooh !!"

	case "remove":
		if len(args) != 2 { return s.Help() }
		vc.Character.SpellsMemorized.Remove(args[1])
		return "it's probably gone"
		
	case "cast":
		if len(args) < 2 { return s.Help() }
		level, name := 0, args[1]
		if len(args) == 3 {
			if v,err := strconv.Atoi(args[2]); err != nil {
				return s.Help() + " " + err.Error()
			} else {
				level = v
			}
		}
		str,_ := vc.Character.CastSpell(name, level)
		return str

	case "refresh":
		vc.Character.Slots.Reset()
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
