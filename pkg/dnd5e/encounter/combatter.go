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
	GetDamagerNames() []string
	GetArmorClass() int
	GetHP() (int, int)
	GetAttr(character.AttrKind) int // Lookup Str,Dex, etc

	TakeDamage(d int)
}

// Damager is something that can be used to make an attack, via a roll
type Damager interface {
	GetName() string
	GetHitModifier() int        // What gets added to the base attack roll
	GetDamageRoll() string      // e.g. "4d6+3". If "", then no damage roll.
}

func CombatterToString(c Combatter) string {
	hp,maxhp := c.GetHP()

	str := fmt.Sprintf("[%s] ", c.GetName())

	if hp <= 0 {
		return str + "is *DEAD*"
	}

	str += fmt.Sprintf("HP:%d/%d, AC:%d", hp, maxhp, c.GetArmorClass())

	for _,name := range c.GetDamagerNames() {
		d := c.GetDamager(name)
		str += fmt.Sprintf(" %s", DamagerToString(d))
	}

	return str
}

func DamagerToString(d Damager) string {
	str := fmt.Sprintf("(%s: ", d.GetName())

	if mod := d.GetHitModifier(); mod != 0 {
		str += fmt.Sprintf("hit:%+d", mod)
	}
	if d.GetDamageRoll() != "" {
		str += ", damage:" + d.GetDamageRoll()
	}

	return str + ")"
}

func CombatterRollInitiative(c Combatter) (int, string) {
	dex := c.GetAttr(character.Dex)
	mod := character.AttrModifier(dex)
	r := roll.Roll{NumDice:1, DiceSize:20, Modifier:mod, Reason:fmt.Sprintf("initiative, dex=%d", dex)}
	o := r.Do()
	return o.Total, o.String()
}
