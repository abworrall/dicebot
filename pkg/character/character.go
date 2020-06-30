package character

import(
	"fmt"
	"strconv"

	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
	"github.com/abworrall/dicebot/pkg/dnd5e/spells"
)

// A Character holds info about a typical RPG character
type Character struct {
	Name string
	Race string
	Class string
	Level int
	Alignment string
	Str,Int,Wis,Con,Cha,Dex,Per int
	MaxHitpoints int
	CurrHitpoints int

	Weapons map[string]int
	CurrWeapon string
	Armor string
	Shield bool

	SpellsMemorized spells.Set
	Slots spells.Slots

	Inventory
}

func NewCharacter() Character {
	return Character{
		Weapons: map[string]int{},
		SpellsMemorized: spells.NewSet(),
		Inventory: NewInventory(),
	}
}

func (c Character)String() string {
	s := fmt.Sprintf(`--{ %s }--
Race: %s
Class: %s (%d)
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
		c.Name, c.Race, c.Class, c.Level, c.Alignment,
		c.Str, AttrModifier(c.Str),
		c.Int, AttrModifier(c.Int),
		c.Wis, AttrModifier(c.Wis),
		c.Con, AttrModifier(c.Con),
		c.Cha, AttrModifier(c.Cha),
		c.Dex, AttrModifier(c.Dex),
		c.Per, AttrModifier(c.Per),
		c.CurrHitpoints, c.MaxHitpoints)

	if c.Armor != "" {
		_,desc := c.GetArmorClass()
		s += "\nArmor: " + desc + "\n"
	}

	if len(c.Weapons) > 0 {
		s += "\n--{ Weapons }--\n"
		for name,_ := range c.Weapons {
			w := rules.TheRules.EquipmentList[name]

			prefix := "   "
			if c.CurrWeapon == name { prefix = "** " }
			hitMod,hitModDesc := c.GetWeaponAttackModifier(w)
			s += fmt.Sprintf("%s%s, hit:%+d %s\n", prefix, w.WeaponDamageString(), hitMod, hitModDesc)
		}
	}

	if c.IsSpellCaster() {
		s += c.MagicString()
	}

	
	
	return s
}

func (c *Character)IsSpellCaster() bool { return c.Slots.Max[1] > 0 }


func (c *Character)MagicString() string {
	if !c.IsSpellCaster() {
		return "You can't do magic :("
	}

	_,desc := c.GetMagicAttackModifier()

	return fmt.Sprintf("\nSpell Attack Modifier: %s\n\n%s\n--- %s", desc, c.SpellsMemorized, c.Slots)
}


func (c *Character)Set(k,v string) string {
	myatoi := func(s string) int {
		i,_ := strconv.Atoi(s)
		return i
	}

	switch k {
	case "name": c.Name = v
	case "race": c.Race = v
	case "class": c.Class = v
	case "alignment": c.Alignment = v

	case "weapon":
		if ! rules.TheRules.IsWeapon(v) {
			return fmt.Sprintf("'%s' is not a known weapon", v)
		}
		c.Weapons[v] = 1
		c.CurrWeapon = v

	case "shield": c.Shield = (myatoi(v) == 1)
		
	case "armor":
		if ! rules.TheRules.IsArmor(v) {
			return fmt.Sprintf("'%s' is not a known kind of armor", v)
		}
		c.Armor = v
		
	case "str": c.Str = myatoi(v)
	case "int": c.Int = myatoi(v)
	case "wis": c.Wis = myatoi(v)
	case "con": c.Con = myatoi(v)
	case "cha": c.Cha = myatoi(v)
	case "dex": c.Dex = myatoi(v)
	case "per": c.Per = myatoi(v)

	case "level": c.Level = myatoi(v)
	case "maxhp": c.MaxHitpoints = myatoi(v)
	case "hp": c.CurrHitpoints = myatoi(v)

	default: return fmt.Sprintf("I don't set '%s'", k)
	}

	return fmt.Sprintf("%s set to %s", k, v)
}

// Get is a gruesome kind of thing.
func (c *Character)Get(k string) (int,string) {
	i := -1
	str := ""

	switch k {
	case "name": str = c.Name
	case "race": str = c.Race
	case "class": str = c.Class
	case "alignment": str = c.Alignment

	case "str": i = c.Str
	case "int": i = c.Int
	case "wis": i = c.Wis
	case "con": i = c.Con
	case "cha": i = c.Cha
	case "dex": i = c.Dex
	case "per": i = c.Per

	case "level": i = c.Level
	case "maxhp": i = c.MaxHitpoints
	case "hp": i = c.CurrHitpoints
	}

	return i,str
}



// Summary returns a oneliner
func (c *Character)Summary() string {
	str := fmt.Sprintf("[%s] L%d %s, HP:%d/%d", c.Name, c.Level, c.Class, c.CurrHitpoints, c.MaxHitpoints)

	if c.Slots.Max[1] != 0 {
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
