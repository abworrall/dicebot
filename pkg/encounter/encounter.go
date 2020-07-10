package encounter

import(
	"fmt"
	"sort"

	"github.com/abworrall/dicebot/pkg/roll"
	"github.com/abworrall/dicebot/pkg/character"
)

type Encounter struct {
	Name string
	Combatants map[string]Combatter
	Init Initiative
	TurnNumber int

	GroupCounts map[string]int   // e.g. GroupCounts["goblin"] = 8 when there are 8 goblins in play
}
func (e Encounter)IsNil() bool { return len(e.Combatants) == 0 }

func NewEncounter() Encounter {
	return Encounter{
		Combatants: map[string]Combatter{},
		Init: NewInitiative(),
		GroupCounts: map[string]int{},
	}
}

func (e Encounter)String() string {
	if e.IsNil() { return "nah, everyone is chillin'" }

	str := ""

	names := []string{}
	for _,v := range e.Combatants {
		names = append(names, v.GetName())
	}
	sort.Strings(names)
	for _,name := range names {
		str += CombatterToString(e.Combatants[name]) + "\n"
	}

	str += "\n" + e.Init.String()

	return str
}

// Returns 1 for monster instance of that name, then 2, etc
func (e *Encounter)NextGroupIndex(name string) int {
	e.GroupCounts[name]++ // auto-inits
	return e.GroupCounts[name]
}

func (e *Encounter)Add(c Combatter) {
	e.Combatants[c.GetName()] = c
}

func (e *Encounter)Lookup(name string) (Combatter, bool) {
	c,exists := e.Combatants[name]
	return c,exists
}

type AttackSpec struct {
	Targets []Combatter

	// For a fully auto attack using the combatter's damager
	Attacker Combatter
	DamagerName string
	WithAdvantage bool
	WithDisadvantage bool
}
func NewAttackSpec() AttackSpec {
	return AttackSpec{Targets: []Combatter{}}
}

func (e *Encounter)Attack(spec AttackSpec) string {	//c1,c2 Combatter, weaponOrAction string) string {
	str := ""

	c1 := spec.Attacker
	weapon := c1.GetDamager(spec.DamagerName)
	if weapon == nil || weapon.GetName() == "" {
		return fmt.Sprintf("%s did not have a '%s'", c1.GetName(), spec.DamagerName)
	}

	for _,c2 := range spec.Targets {
		str += fmt.Sprintf("%s attacks %s with %s: ", c1.GetName(), c2.GetName(), weapon.GetName())

		// Build the hit roll
		hitRoll := roll.Roll{
			NumDice: 1,
			DiceSize: 20,
			Modifier: weapon.GetHitModifier(),
			Target: c2.GetArmorClass(),
			WithAdvantage: spec.WithAdvantage,
			WithDisadvantage: spec.WithDisadvantage,
			WithImprovedCritical: c1.HasBuff(character.BuffFighterChampionImprovedCritical),
		}

		hitOutcome := hitRoll.Do()
		str += "Attack - " + hitOutcome.String()

		if hitOutcome.Success {
			// We hit ! Damage roll ?
			if weapon.GetDamageRoll() != "" {
				damageRoll := roll.Parse(weapon.GetDamageRoll())

				if hitOutcome.CriticalHit {
					// Twice as many damage dice !
					str += " DOUBLE DICE"
					damageRoll.NumDice *= 2
				}

				damageOutcome := damageRoll.Do()
				str += " Damage - " + damageOutcome.String()

				c2.AdjustHP(-1 * damageOutcome.Total)

				hp,_ := c2.GetHP()
				if hp == 0 {
					str += " they are DEAD"
				}
			}
		}
		str += "\n"
	}

	return str
}
