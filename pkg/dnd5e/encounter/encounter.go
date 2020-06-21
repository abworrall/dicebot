package encounter

import(
	"github.com/abworrall/dicebot/pkg/dnd5e/roll"
)

// This is a sketch for now

type BaseAttr int
const(
	Str BaseAttr = iota
	Dex
)

// Weaponer is something that can be used to make an attack, via a roll
type Weaponer interface {
	GetHitModifier() int // What gets added to the base attack role
	GetDamageRoll() string
	GetBaseAttr() BaseAttr
}

// Combatanter represents an entity that will be fighting.
// This will be a wrapped monster, or a wrapped character.
type Combatter interface {
	GetName() string
	GetCurrentWeapon() Weaponer
	GetArmorClass() int
	GetHP() (int, int)
	GetAttr(BaseAttr) int // Lookup the base attribute

	TakeDamage(d int)
}

type Encounter struct {
	Combatants []Combatter     // each is a wrapped clone of a PC or monster
	Initiative map[string]int
	Turn int
}

func (e *Encounter)Add(c Combatter) {}
func (e *Encounter)Attack(c1,c2 Combatter) {
	w := c1.GetCurrentWeapon()

	// Build up the modifiers for the attack
	modifier := w.GetHitModifier() // The weapon might have a +2 or whatever
	modifier += roll.AttrToModifier(c1.GetAttr(w.GetBaseAttr())) // E.g. if STR weapon, get the STR modifier

	// Build the hit roll !
	hit := roll.Roll{
		NumDice: 1,
		DiceSize: 20,
		Modifier: modifier,
		Target: c2.GetArmorClass(),
	}

	// TODO: figure out how to wedge advantage/disadvantage into this
	o := hit.Do()
	if !o.Success {
		// Attack failed
		return
	}

	// Damage roll !
}
