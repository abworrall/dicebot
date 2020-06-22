package character

import(
	"fmt"
	"regexp"
	"strconv"

	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)

// SpellSlots represents the spells the player has memorized, and whether they are ready to cast
type SpellSlots struct {
	Kind      string   // What kind of spells are allowed
	Max     []int      // Max number of spells at each level (index 0 == level 1)
	Memo  [][]Slot    // What's actually memorized in that slot (index 0 == level 1)
}
type Slot struct {
	Name      string
	Spent     bool     // Whether it is ready to use or has already been used
}

// How a slot is identified
type SlotIndex struct {
	Level int
	Index int
}
func (si SlotIndex)String() string { return fmt.Sprintf("%d.%d", si.Level, si.Index) }

// {{{ s.String

func (s Slot)String() string {
	if s.Name == "" {
		return "[empty]"
	}

	sp := rules.TheRules.SpellList.LookupFirst(s.Name)	
	str := sp.Summary()
	
	if s.Spent {
		str = "(*spent*) " + str
	}		

	return str
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
	
	str := fmt.Sprintf("Spell Slots:  (Max:%v, Kind:%s)\n", sl.Max, sl.Kind)
	for lvl, slots := range sl.Memo {
		for idx, slot := range slots {
			str += fmt.Sprintf(" L%d.%d  %s\n", lvl+1, idx+1, slot)
		}
	}

	return str
}

// }}}
// {{{ sl.Init

// SetMax specifies the max slots for each level.
func (sl *SpellSlots)Init(max []int, kind string) {
	sl.Max = max
	sl.Kind = kind
	
	// Presize all the slot slices to the max size
	sl.Memo = make([][]Slot, len(max))
	for i,n := range max {
		sl.Memo[i] = make([]Slot, n)
	}
}

// }}}
// {{{ sl.Memorize

func (sl *SpellSlots)Memorize(sb *Spellbook, idxStr string, name string) error {
	if ! rules.TheRules.IsSpell(name) {
		return fmt.Errorf("'%s' not a known spell (look up name with `db rules spell blah`)", name)
	}

	if sb!=nil && !sb.IsKnown(name) {
		return fmt.Errorf("spell '%s' not in your spellbook", name)
	}
	
	if slot,_,err := sl.Lookup(idxStr); err != nil {
		return err
	} else {
		slot.Name = name
		slot.Spent = false
	}

	return nil
}

// }}}
// {{{ sl.Cast

func (sl *SpellSlots)Cast(idxStr string) (string, error) {
	slot,_,err := sl.Lookup(idxStr)
	if err != nil { return "", err }

	if slot.Spent {
		return "", fmt.Errorf("That spell slot is spent :(")
	}

	slot.Spent = true
	return slot.Name, nil
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
func (sl *SpellSlots)Lookup(idxStr string) (*Slot, SlotIndex, error) {
	si,err := parseSlotIndex(idxStr)
	if err != nil { return nil, SlotIndex{}, err }

	if si.Level > len(sl.Memo) {
		return nil, SlotIndex{}, fmt.Errorf("You don't know things as fancy as level %d", si.Level)
	}
	
	// SlotIndex is 1-indexed; since we're going to index into arrays, decrement it
	lvl := si.Level - 1
	idx := si.Index - 1
	
	if idx >= len(sl.Memo[lvl]) {
		return nil, SlotIndex{}, fmt.Errorf("You only get %d spells at level %d", len(sl.Memo[lvl]), si.Level)
	} else {
		return &sl.Memo[lvl][idx], si, nil
	}
}

// }}}

// {{{ parseSlotIndex

// ParseIndex parse an index style string (e.g. `2:4`) into (level, index), or returns an error.
// Note that the retrurned level value starts at one, but the returned index value starts at zero (for slice lookup)
func parseSlotIndex(s string) (SlotIndex, error) {
	bits := regexp.MustCompile(`^L?(\d*)[.:](\d+)$`).FindStringSubmatch(s) // 2.4, 2, 4
	if len(bits) != 3 {
		return SlotIndex{}, fmt.Errorf("index `%s` is nonsense", s)
	}
	lvl,_   := strconv.Atoi(bits[1])
	idx,_ := strconv.Atoi(bits[2])

	if lvl < 1 || lvl > 9  { return SlotIndex{}, fmt.Errorf("you want level %d ? you can't handle level %d", lvl, lvl) }
	if idx < 1 || idx > 15 { return SlotIndex{}, fmt.Errorf("index %d is mad index", idx) }
	return SlotIndex{Level:lvl, Index:idx}, nil
}

// }}}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
