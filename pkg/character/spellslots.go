package character

import(
	"fmt"

	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)

// SpellSlots represents the spells the player has memorized, and whether they are ready to cast
type SpellSlots struct {
	Max     []int      // Max number of spells at each level (index 0 == level 1)
	Memo  [][]Slot     // What's actually memorized in that slot (index 0 == level 1)
}
type Slot struct {
	Spell    *Spell    // Points into the Spell in the Spellbook
	Spent     bool     // Whether it is ready to use or has already been used
}

type SlotIndex string // Should be of the form `2:4`, e.g. the fourth 2nd-level slot

// {{{ s.String

func (s Slot)String() string {
	if s.Spell == nil {
		return "[empty]"
	}

	descrip := "UNRECOGNIZED: "+s.Spell.String()

	if descrips := rules.TheRules.SpellList.Lookup(s.Spell.Name); len(descrips) == 1 {
		descrip = descrips[0].Summary()
	}
	
	if s.Spent {
		return fmt.Sprintf("(%s)", descrip)
	}		

	return descrip
}

// }}}

// {{{ NewSpellSlots

func NewSpellSlots() SpellSlots {
	return SpellSlots{
		Max:     []int{},
		Memo:  [][]Slot{},
	}
}

// }}}
// {{{ sl.String

func (sl *SpellSlots)String() string {
	if sl == nil { return "" }

	if len(sl.Memo) == 0 { return "You can't do magic spells :(" }
	
	str := fmt.Sprintf("Spell Slots:  (Max:%v)\n", sl.Max)
	for lvl, slots := range sl.Memo {
		for idx, slot := range slots {
			str += fmt.Sprintf(" L%d:%d  %s\n", lvl+1, idx+1, slot)
		}
	}

	return str
}

// }}}
// {{{ sl.SetMax

// SetMax specifies the max slots for each level.
func (sl *SpellSlots)SetMax(max []int) {
	sl.Max = max

	// Presize all the slot slices to the max size
	sl.Memo = make([][]Slot, len(max))
	for i,n := range max {
		sl.Memo[i] = make([]Slot, n)
	}
}

// }}}
// {{{ sl.Memorize

// Memorize looks up spell `i` from the spellbook, and stores it in
// the slots, as per the level and `idx`. If there is a problem,
// returns an error.
func (sl *SpellSlots)Memorize(sb *Spellbook, i SpellIndex, idx int) error {
	spell,err := sb.Lookup(i)
	if err != nil { return err }

	idxStr := SlotIndex(fmt.Sprintf("%d:%d", spell.Level, idx))

	return sl.MemorizeSpell(spell, idxStr)
}

// }}}
// {{{ sl.MemorizeSpell

// Memorize looks up spell `i` from the spellbook, and stores it in
// the slots, as per the level and `idx`. If there is a problem,
// returns an error.
func (sl *SpellSlots)MemorizeSpell(spell *Spell, idx SlotIndex) error {
	if slot,err := sl.Lookup(idx); err != nil {
		return err
	} else {
		slot.Spell = spell
		slot.Spent = false
	}

	return nil
}

// }}}
// {{{ sl.Cast

func (sl *SpellSlots)Cast(i SlotIndex) (*Spell, error) {
	slot,err := sl.Lookup(i)
	if err != nil { return nil, err }

	if slot.Spent {
		return nil, fmt.Errorf("That spell slot is spent :(")
	}

	slot.Spent = true
	return slot.Spell, nil
}

// }}}
// {{{ sl.Refresh

// Refresh rememorizes all the spells, so they are ready for use
func (sl *SpellSlots)Refresh() {
	for lvl, slots := range sl.Memo {
		for idx, _ := range slots {
			sl.Memo[lvl][idx].Spent = false
		}
	}
}

// }}}

// {{{ sl.Lookup

// Lookup will lookup the slot, or return an error
func (sl *SpellSlots)Lookup(i SlotIndex) (*Slot, error) {
	lvl,idx,err := ParseIndex(string(i))
	if err != nil { return nil, err }

	if lvl > len(sl.Memo) {
		return nil, fmt.Errorf("You don't know things as fancy as level %d", lvl)
	}
	
	// lvl is 1-indexed; since we're going to index into arrays, decrement it
	lvl = lvl - 1

	if idx > len(sl.Memo[lvl]) {
		return nil, fmt.Errorf("You only get %d spells at level %d", len(sl.Memo[lvl]), lvl+1)
	} else {
		return &sl.Memo[lvl][idx], nil
	}
}

// }}}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
