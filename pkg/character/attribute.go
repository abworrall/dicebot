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

func AttrModifier(attrVal int) int {
	// Kinda horrible lookup table
	var attrModifiers = []int{
		0,-5,-4,-4,          // attr scores 0-3
		-3,-3,-2,-2,-1,-1,   // attr scores 4-9
		0,0,1,1,2,2,         // attr scores 10-15
		3,3,4,4,5,5,         // attr scores 16-21
		6,6,7,7,8,8,9,9,10,  // attr scores 22-30
	}
	if attrVal<0 || attrVal>len(attrModifiers) {
		return 0
	}
	return attrModifiers[attrVal]
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
