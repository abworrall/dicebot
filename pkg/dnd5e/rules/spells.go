package rules

import(
	"fmt"
	"strings"
)

type Spell struct{
	Index string `json:"index"`
	Name string `json:"name"`
	Desc[] string `json:"desc"`
	Higher[] string `json:"higher_level"`
	Range string `json:"range"`
	Duration string `json:"duration"`
	Level int `json:"level"`
	Classes []struct{
		Name string `json:"name"`
	} `json:"classes"`
}

func (s Spell)String() string { return s.Summary() }

func (s Spell)Type() string { return "spell" }

func (s Spell)Summary() string {
	return fmt.Sprintf("[L%d (%s), %s] %s", s.Level, s.Class(), s.Index, s.Name)
}

func (s Spell)Description() string {
	return fmt.Sprintf(`--{ %s }--
Level: %d (%s)
Range: %s
Duration: %s
%s
%s`, 
	s.Name,
	s.Level, s.Class(),
	s.Range, 
	s.Duration, 
	s.Desc,
	s.Higher)
}


// SpellList just maps the `Index` of each spell to the spell object
type SpellList map[string]Spell

// Implement the lookup interface
func (sl SpellList)Lookup(namelike string) []Entryer {
	if v,exists := sl[namelike]; exists {
		return []Entryer{v}
	}

	ret := []Entryer{}
	for _,v := range sl {
		if strings.Contains(strings.ToLower(v.Name),strings.ToLower(namelike)) {
			ret = append(ret, Entryer(v))
		}
	}
	return ret
}

func (s Spell)Class() string {
	names := []string{}
	for _,c := range s.Classes {
		names = append (names, c.Name)
	}
	return strings.Join(names, ",")
}
