package verbs

import(
	"fmt"
	"strings"

	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/roll"
)

type Roll struct{}

func (r Roll)Help() string {
	return "[4d6+3 >=8 withadvantage] [init]"
}

func (r Roll)Process(vc VerbContext, args []string) string {
	if len(args) < 1 {
		return "what do you want me to roll ?"
	}

	// Some named rolls
	switch args[0] {
	case "init": return r.RollInitiative(vc)
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

	dex := vc.Character.GetAttr(character.Dex)
	mod := character.AttrModifier(dex)
	reason := fmt.Sprintf("initiative for %s, dex=%d", name, dex)
	initRoll := roll.Roll{NumDice:1, DiceSize:20, Modifier:mod, Reason:reason}
	o := initRoll.Do()

	vc.Encounter.Init.Set(name, o.Total)

	return o.String()
}
