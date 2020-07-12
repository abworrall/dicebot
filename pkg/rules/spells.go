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
	CastingTime string `json:"casting_time"`
	Ritual bool `json:"ritual"`
	Level int `json:"level"`
	Classes []struct{
		Name string `json:"name"`
	} `json:"classes"`
	Subclasses []struct{
		Name string `json:"name"`
	} `json:"subclasses"`
}
func (s Spell)IsNil() bool { return s.Index == "" }

func (s Spell)String() string { return s.Summary() }

func (s Spell)Type() string { return "spell" }

func (s Spell)Summary() string {
	return fmt.Sprintf("[L%d (%s/%s), %s, cast:{%s}]", s.Level, s.ClassesString(), s.SubclassesString(), s.Index, s.CastingTime)
}

func (s Spell)ShorterSummary() string {
	cast := s.CastingTime
	if s.Ritual {
		cast = "ritual, " + cast
	}
	return fmt.Sprintf("L%d %s cast:{%s}", s.Level, s.Index, cast)
}

func (s Spell)Description() string {
	return fmt.Sprintf(`--{ %s }--
Level: %d (%s) (%s)
Range: %s
Duration: %s
Casting time: %s
Ritual: %v
%s
%s`, 
	s.Name,
	s.Level, s.ClassesString(), s.SubclassesString(),
	s.Range, 
	s.Duration,
	s.CastingTime,
	s.Ritual,
	s.Desc,
	s.Higher)
}

func (s Spell)HasClass(c string) bool {
	for _,v := range s.Classes {
		if strings.EqualFold(v.Name, c) {
			return true
		}
	}
	return false
}

func (s Spell)HasSubclass(c string) bool {
	for _,sc := range s.Subclasses {
		if strings.EqualFold(sc.Name, c) {
			return true
		}
	}
	return false
}

func (s Spell)ClassesString() string {
	names := []string{}
	for _,c := range s.Classes {
		names = append (names, c.Name)
	}
	return strings.Join(names, ",")
}

func (s Spell)SubclassesString() string {
	names := []string{}
	for _,c := range s.Subclasses {
		names = append (names, c.Name)
	}
	return strings.Join(names, ",")
}


// SpellList just maps the `Index` of each spell to the spell object
type SpellList map[string]Spell


func (sl SpellList)FindMatching(class, subclass string, level int) []Spell {
	ret := []Spell{}

	for _,spell := range sl {
		if level > 0 && spell.Level > level { continue }
		if class != "" && !spell.HasClass(class) { continue }
		if subclass != "" && !spell.HasSubclass(subclass) { continue }
		ret = append(ret, spell)
	}
	
	return ret
}

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

func (sl SpellList)LookupFirst(namelike string) Entryer {
	if m := sl.Lookup(namelike); len(m) > 0 {
		return m[0]
	}
	return nil
}
