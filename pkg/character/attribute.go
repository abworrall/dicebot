package character

type AttrKind int
const(
	Str AttrKind = iota
	Int
	Wis
	Con
	Cha
	Dex
	Per // Remove sometime ?
)

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
