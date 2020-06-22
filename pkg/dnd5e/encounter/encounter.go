package encounter

import(
	"fmt"
	"sort"
	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/roll"
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
	if e.IsNil() { return "no data" }

	str := fmt.Sprintf("--{ %d participants }--\n", len(e.Combatants))

	if len(e.Combatants) == 0 {
		return str
	}

	str += e.Init.String()+"\n"

	names := []string{}
	for _,v := range e.Combatants {
		names = append(names, v.GetName())
	}
	sort.Strings(names)
	for _,name := range names {
		str += CombatterToString(e.Combatants[name]) + "\n"
	}

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

func (e *Encounter)Attack(c1,c2 Combatter, weaponOrAction string) {

	w := c1.GetDamager(weaponOrAction)

	// Build up the modifiers for the attack
	modifier := w.GetHitModifier() // The weapon might have a +2 or whatever
	modifier += character.AttrModifier(c1.GetAttr(w.GetModifierAttr())) // E.g. if STR weapon, get the STR modifier

	// Build the hit roll
	hit := roll.Roll{
		NumDice: 1,
		DiceSize: 20,
		Modifier: modifier,
		Target: c2.GetArmorClass(),
		// TODO: figure out how to wedge advantage/disadvantage into this
	}

	hitOutcome := hit.Do()
	if !hitOutcome.Success {
		// Attack failed
		return
	}

	// Damage roll !
	damage := roll.Parse(w.GetBaseDamageRoll())
	// TODO: figure if this is double counting for monsters, who maybe shouldn't get attr based modifiers as well
	damage.Modifier += w.GetDamageModifier()
	damage.Modifier += character.AttrModifier(c1.GetAttr(w.GetModifierAttr()))

	damageOutcome := damage.Do()

	c2.TakeDamage(damageOutcome.Total)
}
