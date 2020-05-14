package character

import(
	"fmt"
	"strconv"
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

	Inventory
}

func NewCharacter() Character {
	return Character{
		Inventory: Inventory{
			Items: []Item{},
		},
	}
}

func (c Character)String() string {
	s := fmt.Sprintf(`--{ %s }--
Race: %s
Class: %s (%d)
Alignment: %s

STR: %2d
INT: %2d
WIS: %2d
CON: %2d
CHA: %2d
DEX: %2d
PER: %2d

HP: (%d/%d)
`,
		c.Name, c.Race, c.Class, c.Level, c.Alignment,
		c.Str,
		c.Int,
		c.Wis,
		c.Con,
		c.Cha,
		c.Dex,
		c.Per,
		c.CurrHitpoints, c.MaxHitpoints)
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