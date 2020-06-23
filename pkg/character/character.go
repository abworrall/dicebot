package character

import(
	"fmt"
	"strconv"

	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
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
	
	Inventory
	Spellbook
	SpellSlots
}

func NewCharacter() Character {
	return Character{
		Inventory: NewInventory(),
		Spellbook: NewSpellbook(),
		SpellSlots: NewSpellSlots(),
		Weapons: map[string]int{},
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
		armor := rules.TheRules.EquipmentList.LookupFirst(c.Armor)
		s += fmt.Sprintf("\nArmor: %s", armor.Summary())
		if c.Shield {
			s += ", and a shield too"
		}
		s += "\n"
	}

	if len(c.Weapons) > 0 {
		s += "\n--{ Weapons }--\n"
		for name,_ := range c.Weapons {
			w := rules.TheRules.EquipmentList.LookupFirst(name)

			prefix := "   "
			if c.CurrWeapon == name { prefix = "** " }
			s += fmt.Sprintf("%s%s\n", prefix, w.Summary())
		}
	}

	if len(c.SpellSlots.Memo) > 0 {
		s += fmt.Sprintf("\n%s", c.SpellSlots.String())
	}
	
	return s
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
