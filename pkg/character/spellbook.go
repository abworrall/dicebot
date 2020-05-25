package character

import(
	"fmt"
	"regexp"
	"strconv"
)

// A Spellbook is a list of all the spells that the character knows, organized by level.
type Spellbook struct {
	Spells map[int][]Spell  // map key is level
}

func NewSpellbook() Spellbook {
	return Spellbook{
		Spells: map[int][]Spell{},
	}
}

// A Spell is a single castable thing. Perhaps with effects, descriptions, etc.
type Spell struct {
	Level int
	Name string
}
func (sp Spell)String() string { return sp.Name }

type SpellIndex string // Should be of the form `2:4`, e.g. the fourth 2nd-level slot

// {{{ sb.String

func (sb *Spellbook)String() string {
	if sb == nil || sb.Spells == nil { return "" }
	
	str := "Spellbook:-\n"
	for lvl := 1; lvl<20; lvl++ {
		if spells,exists := sb.Spells[lvl]; !exists {
			continue
		} else {
			for i,s := range spells {
				str += fmt.Sprintf(" L%d:%d %s\n", s.Level, i+1, s)
			}
		}
	}
	return str
}

// }}}
// {{{ sb.Learn

// Learn adds a new spell into the spellbook. Does not check for dupes.
func (sb *Spellbook)Learn(level int, name string) {
	if _,exists := sb.Spells[level]; !exists {
		sb.Spells[level] = []Spell{}
	}

	sb.Spells[level] = append(sb.Spells[level], Spell{level, name})
}

// }}}
// {{{ sb.Lookup

// Lookup finds the spell specified by the index, or returns an error.
func (sb *Spellbook)Lookup(si SpellIndex) (*Spell, error) {
	if lvl,idx,err := ParseIndex(string(si)); err != nil {
		return nil, err
	} else if _,exists := sb.Spells[lvl]; !exists {
		return nil, fmt.Errorf("You don't know *any* level %d spells", lvl)
	} else if idx > len(sb.Spells[lvl]) {
		return nil, fmt.Errorf("You only know %d level %d spells", len(sb.Spells[lvl]), lvl)
	} else {
		return &sb.Spells[lvl][idx], nil
	}
}

// }}}

// {{{ ParseIndex

// ParseIndex parse an index style string (e.g. `2:4`) into (level, index), or returns an error.
// Note that the retrurned level value starts at one, but the returned index value starts at zero (for slice lookup)
func ParseIndex(s string) (int, int, error) {
	bits := regexp.MustCompile(`^(\d*):(\d+)$`).FindStringSubmatch(s) // 2:4, 2, 4
	if len(bits) != 3 {
		return 0,0,fmt.Errorf("index `%s` is nonsense", s)
	}
	lvl,_   := strconv.Atoi(bits[1])
	idx,_ := strconv.Atoi(bits[2])

	if lvl < 1 || lvl > 9  { return 0,0,fmt.Errorf("you want level %d ? you can't handle level %d", lvl, lvl) }
	if idx < 1 || idx > 15 { return 0,0,fmt.Errorf("index %d is mad index", idx) }
	return lvl, idx-1, nil
}

// }}}


// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
