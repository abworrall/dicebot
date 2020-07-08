package character

import	"strings"

type AttrKind int
const(
	Undef AttrKind = iota
	Str
	Int
	Wis
	Con
	Cha
	Dex
	Per // Remove sometime ?
)

func (a AttrKind)String() string {
	switch a {
	case Str: return "str"
	case Int: return "int"
	case Wis: return "wis"
	case Con: return "con"
	case Cha: return "cha"
	case Dex: return "dex"
	case Per: return "per"
	default:  return "???"
	}
}

func ParseAttr(s string) AttrKind {
	switch strings.ToLower(s) {
	case "str": return Str
	case "int": return Int
	case "wis": return Wis
	case "con": return Con
	case "cha": return Cha
	case "dex": return Dex
	case "per": return Per
	default:    return Undef
	}
}

func (c Character)GetAttr(k AttrKind) int {
	switch k {
	case Str: return c.Str
	case Int: return c.Int
	case Wis: return c.Wis
	case Con: return c.Con
	case Cha: return c.Cha
	case Dex: return c.Dex
	case Per: return c.Per
	default:  return -1
	}
}

func (c Character)GetModifier(k AttrKind) int {
	return AttrModifier(c.GetAttr(k))
}

func (c Character)GetAttrAndModifier(k AttrKind) (int,int) {
	return c.GetAttr(k), AttrModifier(c.GetAttr(k))
}


func (c *Character)GetSpellcastingAbilityAttr() AttrKind {
	switch c.Class {
	case "wizard":
		return Int

	case "ranger": fallthrough
	case "cleric": fallthrough
	case "druid":
		return Wis

	case "bard": fallthrough
	case "paladin": fallthrough
	case "warlock": fallthrough
	case "sorceror":
		return Cha
		
	default:
		return Undef
	}
}
