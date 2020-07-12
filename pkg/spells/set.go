package spells

import(
	"fmt"
	"sort"
)

// A Set is a list of spells with no dupes. It is the basis for a
// spellbook, and also for any magic-user's list of prepared spells.
type Set struct {
	Kind string
	Max int
	Spells map[string]int  // A val of 1 means hand-selected by user; 2 means from other sources
}

func NewSet() Set {
	return Set{
		Spells: map[string]int{},
	}
}

func (s Set)String() string {
	if s.Spells == nil { return "" }

	str := fmt.Sprintf("-- Spellset (max:%d, kind:%s)\n", s.Max, s.Kind)
	str2 := ""

	keys := []string{}
	for name,_ := range s.Spells {
		keys = append(keys, name)
	}

	// Custom sort function
	sort.Slice(keys, func(i, j int) bool {
		s1,s2 := Lookup(keys[i]), Lookup(keys[j])
		if s1.Level != s2.Level {
			return s1.Level < s2.Level
		}
		return s1.Index < s2.Index
	})

	for _,name := range keys {
		val := s.Spells[name]
		if val == 2 {
			str2 += "{" + Lookup(name).ShorterSummary() + "}\n"
		} else {
			str += Lookup(name).ShorterSummary() + "\n"
		}
	}
	return str + str2
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

// Clone copies the set into an entirely separate data structure
func (s Set)Clone() Set {
	new := NewSet()
	new.Max = s.Max
	new.Kind = s.Kind
	for k,v := range s.Spells {
		new.Spells[k] = v
	}
	return new
}

// Merge includes one set into another, tagging new ones with the integer val
func (s1 Set)UnionWith(s2 Set, i int) Set {
	new := s1.Clone()
	for k,_ := range s2.Spells {
		new.Spells[k] = i
	}
	return new
}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
