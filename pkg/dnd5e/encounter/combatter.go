package encounter

import(
	"fmt"
	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/roll"
)

// Combatter represents an entity that will be fighting. This will be
// a wrapped monster, or a wrapped character.
type Combatter interface {
	GetName() string
	GetGroup() string // for monster, this is just its type, e.g. "goblin"
	GetDamager(name string) Damager
	GetArmorClass() int
	GetHP() (int, int)
	GetAttr(character.AttrKind) int // Lookup Str,Dex, etc

	TakeDamage(d int)
}

// Damager is something that can be used to make an attack, via a roll
type Damager interface {
	GetHitModifier() int        // What gets added to the base attack roll
	GetDamageModifier() int     // What gets added to the base damage roll
	GetBaseDamageRoll() string  // e.g. "4d6+3"
	GetModifierAttr() character.AttrKind      // Which Attr drives the modifier
}

func CombatterToString(c Combatter) string {
	hp,maxhp := c.GetHP()
	return fmt.Sprintf("[%s] HP:%d/%d, AC:%d [Str:%d, Int:%d, Dex:%d]",
		c.GetName(), hp, maxhp, c.GetArmorClass(),
		c.GetAttr(character.Str),
		c.GetAttr(character.Int),
		c.GetAttr(character.Dex))
}

func CombatterRollInitiative(c Combatter) (int, string) {
	dex := c.GetAttr(character.Dex)
	mod := character.AttrModifier(dex)
	r := roll.Roll{NumDice:1, DiceSize:20, Modifier:mod, Reason:fmt.Sprintf("initiative, dex=%d", dex)}
	o := r.Do()
	return o.Total, o.String()
}
