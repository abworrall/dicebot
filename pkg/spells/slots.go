package spells

import(
	"fmt"
	"strings"
)

type Slots struct {
	Max [12]int  // We ignore the element at index 0
	Curr [12]int
}

func NewSlots() Slots {
	return Slots{}
}

func (s Slots)String() string {
	strs := []string{}
	for i:=1; i<len(s.Curr); i++ {
		if s.Max[i] == 0 { continue }
		strs = append(strs, fmt.Sprintf("L%d:%d/%d", i, s.Curr[i], s.Max[i]))
	}
	return fmt.Sprintf("Slots{%s}", strings.Join(strs,", "))
}

func (s *Slots)Spend(level int) bool {
	if level<1 || level>=len(s.Curr) {
		return false
	} else if s.Curr[level] == 0 {
		return false
	}
	s.Curr[level]--
	return true
}

func (s *Slots)Reset() {
	for i:=1; i<len(s.Curr); i++ {
		s.Curr[i] = s.Max[i]
	}
}
