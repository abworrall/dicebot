package verbs

import(
	"fmt"
	"strconv"
)

// HitPoints is stateless, in that the verb doesnb't have its own
// state; it operates on the character's state in the context
type HitPoints struct{}

func (hp HitPoints)Help() string { return "[+n], [-n]" }

func (hp HitPoints)Process(vc VerbContext, args []string) string {
	if vc.Character == nil { return "who are you, again ?" }

	if len(args) != 1 { return hp.Help() }

	if n,err := strconv.Atoi(args[0]); err != nil {
		return fmt.Sprintf("idiocy: %s ", err)
	} else {
		vc.Character.CurrHitpoints += n
	}

	if vc.Character.CurrHitpoints > vc.Character.MaxHitpoints {
		vc.Character.CurrHitpoints = vc.Character.MaxHitpoints
	} else if vc.Character.CurrHitpoints <= 0 {
		return fmt.Sprintf("%s is dead :(", vc.User)
	}

	return fmt.Sprintf("ok, now %d/%d", vc.Character.CurrHitpoints, vc.Character.MaxHitpoints)
}
