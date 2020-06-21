package encounter

import(
	"github.com/abworrall/dicebot/pkg/dnd5e/roll"
)

// This is all a sketch for now, to help figure out the interfaces


type BaseAttr int
const(
	Str BaseAttr = iota
	Dex
)

// Weaponer is something that can be used to make an attack, via a roll
type Weaponer interface {
	GetHitModifier() int        // What gets added to the base attack roll
	GetDamageModifier() int     // What gets added to the base damage roll
	GetBaseDamageRoll() string  // e.g. "4d6+3"
	GetBaseAttr() BaseAttr      // Melee weaponers have Str, or maybe Dex (finesse, blah)
}

// Combatter represents an entity that will be fighting. This will be
// a wrapped monster, or a wrapped character.
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
	TurnNumber int
}

func (e *Encounter)Add(c Combatter) {}

func (e *Encounter)Attack(c1,c2 Combatter) {
	w := c1.GetCurrentWeapon()

	// Build up the modifiers for the attack
	modifier := w.GetHitModifier() // The weapon might have a +2 or whatever
	modifier += roll.AttrToModifier(c1.GetAttr(w.GetBaseAttr())) // E.g. if STR weapon, get the STR modifier

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
	damage.Modifier += roll.AttrToModifier(c1.GetAttr(w.GetBaseAttr()))

	damageOutcome := damage.Do()

	c2.TakeDamage(damageOutcome.Total)
}
