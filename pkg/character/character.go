package character

import(
	"fmt"

	"github.com/abworrall/dicebot/pkg/rules"
	"github.com/abworrall/dicebot/pkg/spells"
)

// A Character holds info about a typical RPG character
type Character struct {
	Name string

	Race string
	Class string
	Subclass string
	Level int
	Alignment string

	Str,Int,Wis,Con,Cha,Dex,Per int

	MaxHitpoints int
	CurrHitpoints int

	Weapons map[string]int

	Armor string	
	CurrWeapon string
	Shield bool

	SpellsMemorized spells.Set
	Slots spells.Slots

	Buffs map[string]int

	Inventory
}

func NewCharacter() Character {
	return Character{
		Weapons: map[string]int{},
		SpellsMemorized: spells.NewSet(),
		Buffs: map[string]int{},
		Inventory: NewInventory(),
	}
}

func (c *Character)IsSpellCaster() bool { return c.Slots.Max[1] > 0 }

func (c *Character)String() string {
	subclass := ""
	if c.Subclass != "" {
		subclass = ", " + c.Subclass
	}

	s := fmt.Sprintf(`--{ %s }--
Race: %s
Class: %s (%d)%s
Alignment: %s

STR: %2d (%+d)
INT: %2d (%+d)
WIS: %2d (%+d)
CON: %2d (%+d)
CHA: %2d (%+d)
DEX: %2d (%+d)
PER: %2d (%+d)

HP: (%d/%d)
`,
		c.Name,
		c.Race,
		c.Class, c.Level, subclass,
		c.Alignment,
		c.Str, AttrModifier(c.Str),
		c.Int, AttrModifier(c.Int),
		c.Wis, AttrModifier(c.Wis),
		c.Con, AttrModifier(c.Con),
		c.Cha, AttrModifier(c.Cha),
		c.Dex, AttrModifier(c.Dex),
		c.Per, AttrModifier(c.Per),
		c.CurrHitpoints, c.MaxHitpoints)

	s += "\n--{ Buffs }--\n"
	for name,_ := range c.Buffs {
		s += name + "\n"
	}
	for name,_ := range c.AutoBuffs() {
		s += "{" + name + "}\n"
	}

	if len(c.Weapons) > 0 {
		s += "\n--{ Weapons"
		if c.CurrWeapon != "" {
			s += " curr=" + c.CurrWeapon
		}
		s += " }--\n"
		for name,_ := range c.Weapons {
			w := rules.TheRules.EquipmentList[name]

			s += fmt.Sprintf("[%s] ", name)

			hitMod,hitModDesc := c.GetWeaponAttackModifier(w)
			damageRoll := c.GetWeaponDamageRoll(w)
			_,damModDesc := c.GetWeaponDamageModifier(w)

			s += fmt.Sprintf("hit:%+d, dam:%s, hit(%s), damage(%s)\n", hitMod, damageRoll, hitModDesc, damModDesc)
		}
	}

	_,desc := c.GetArmorClass()
	s += "\n--{ Armor }--\n" + desc + "\n"

	if c.IsSpellCaster() {
		s += c.MagicString()
	}

	return s
}

func (c *Character)MagicString() string {
	if !c.IsSpellCaster() {
		return "You can't do magic :("
	}
	_,desc := c.GetMagicAttackModifier()
	return fmt.Sprintf("\nSpell Attack Modifier: %s\n\n%s\n--- %s", desc, c.SpellsMemorized, c.Slots)
}

// Summary returns a oneliner
func (c *Character)Summary() string {
	subclass := ""
	if c.Subclass != "" {
		subclass = "{" + c.Subclass + "}"
	}

	str := fmt.Sprintf("[%s] L%d %s%s, HP:%d/%d", c.Name, c.Level, c.Class, subclass, c.CurrHitpoints, c.MaxHitpoints)

	if c.IsSpellCaster() {
		str += ", " + c.Slots.String()
	}
	if c.Armor != "" {
		str += ", " + c.Armor
	}
	if c.CurrWeapon != "" {
		str += ", " + c.CurrWeapon
	}

	return str
}
