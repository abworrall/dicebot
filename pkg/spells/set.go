package spells

import(
	"fmt"
)

// A Set is a list of spells with no dupes. It is the basis for a
// spellbook, and also for any magic-user's list of prepared spells.
type Set struct {
	Kind string
	Max int
	Spells map[string]int
}

func NewSet() Set {
	return Set{
		Spells: map[string]int{},
	}
}

func (s Set)String() string {
	if s.Spells == nil { return "" }

	str := fmt.Sprintf("-- Spellset (max:%d, kind:%s)\n", s.Max, s.Kind)
	for name,_ := range s.Spells {
		str += " " + Lookup(name).ShorterSummary()
	}
	return str
}

// Learn adds a new spell into the set, if it is in the rules
func (s Set)Add(name string) error {
	if s.Max > 0 && len(s.Spells) >= s.Max {
		return fmt.Errorf("Set full, max %d", s.Max)
	}

	if Lookup(name).IsNil() {
		return fmt.Errorf("spell '%s' was not known", name)
	}

	s.Spells[name] = 1
	return nil
}

func (s Set)Remove(name string) {
	delete(s.Spells, name)
}

func (s Set)Contains(name string) bool {
	_,exists := s.Spells[name]
	return exists
}

func (s Set)Size() int {
	return len(s.Spells)
}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
