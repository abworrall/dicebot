package character

import(
	"fmt"

	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)

// A Spellbook is a list of all the spells that the character knows. Key must be a
// recofnized short name (index) as per the rules.
type Spellbook struct {
	Spells map[string]int
}

func NewSpellbook() Spellbook {
	return Spellbook{
		Spells: map[string]int{},
	}
}

// A Spell is a single castable thing.
type Spell struct {
	Level int
	Name string
}
func (sp Spell)String() string { return sp.Name }

type SpellIndex string // Should be of the form `2:4`, e.g. the fourth 2nd-level slot

func (sb *Spellbook)String() string {
	if sb == nil || sb.Spells == nil { return "" }
	
	str := "Spellbook:-\n"
	for name,_ := range sb.Spells {
		str += fmt.Sprintf(" L? %s\n", name)
	}
	return str
}

// Learn adds a new spell into the spellbook. Does not check for dupes.
func (sb *Spellbook)Learn(name string) {
	if rules.TheRules.IsSpell(name) {
		sb.Spells[name] = 1
	}
}

// IsKnown checks to see if the spell is in the book. If the book is
// empty, it assumes the caster knows all the spells (e.g. cleric)
func (sb *Spellbook)IsKnown(name string) bool {
	if len(sb.Spells) == 0 {
		return true
	} else {
		_,exists := sb.Spells[name]
		return exists
	}
}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
