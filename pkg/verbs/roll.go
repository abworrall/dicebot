package verbs

import(
	"fmt"
	"strings"
	"strconv"

	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/roll"
)

type Roll struct{}

func (r Roll)Help() string {
	return "[4d6+3 >=8 withadvantage], [init], [check ATTR [DC] [withadvantage]]"
}

func (r Roll)Process(vc VerbContext, args []string) string {
	if len(args) < 1 {
		return "what do you want me to roll ?"
	}

	// Some named rolls
	switch args[0] {
	case "init": return r.RollInitiative(vc)
	case "check": return r.RollAbilityCheck(vc, args[1:])
	}
	
	str := roll.New(strings.Join(args, " ")).String()
	vc.LogEvent("rolled " + str)
	return str
}

// RollInitiative looks up the current character's Dex, does a roll, and
// then adds them to the current encounter. The various data structures
// are all handled via the VerbContext.
func (r Roll)RollInitiative(vc VerbContext) string {
	if vc.Character == nil || vc.Encounter == nil {
		return "we're not set up for that"
	}

	name := vc.Character.Name
	if val := vc.Encounter.Init.Get(name); val > 0 {
		return fmt.Sprintf("%s has already rolled initiative - it is %d", name, val)
	}

	dex, mod := vc.Character.GetAttrAndModifier(character.Dex)
	reason := fmt.Sprintf("initiative for %s, dex=%d", name, dex)
	initRoll := roll.Roll{NumDice:1, DiceSize:20, Modifier:mod, Reason:reason}
	o := initRoll.Do()

	vc.Encounter.Init.Set(name, o.Total)

	return o.String()
}

// roll check str MM [withadvantage,withdisadvantage]
func (r Roll)RollAbilityCheck(vc VerbContext, args []string) string {
	if vc.Character == nil || len(args) == 0 {
		return "we're not set up for that"
	}

	attrKind := character.ParseAttr(args[0])
	attrVal, attrMod := vc.Character.GetAttrAndModifier(attrKind)

	checkRoll := roll.Roll{
		NumDice: 1,
		DiceSize: 20,
		Modifier: attrMod,
		Reason: fmt.Sprintf("%s ability check %s=%d", vc.Character.Name, attrKind, attrVal),
	}

	args = args[1:]
	
	for len(args) > 0 {
		switch args[0] {
		case "withadvantage": checkRoll.WithAdvantage = true
		case "withdisadvantage": checkRoll.WithDisadvantage = true
		default:
			if n,err := strconv.Atoi(args[0]); err == nil {
				checkRoll.Target = n
			} else {
				return "could not understand this ability check"
			}
		}
		args = args[1:]
	}

	o := checkRoll.Do()

	return o.String()
}
