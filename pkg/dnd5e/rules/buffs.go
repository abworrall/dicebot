package rules

import(
	"fmt"
	"strings"
)

// A Buff is a kinda catchall for a class feature or ability, or some
// other exception to the rules that we want to be aware of.
type Buff struct {
	Index string `json:"index"`

	Class struct{
		Name string `json:"name"`
	} `json:"class"`

	Subclass struct{
		Name string `json:"name"`
	} `json:"subclass"`

	Name string `json:"name"`
	Level int `json:"level"`
	Desc[] string `json:"desc"`
	Group string `json:"group"`
	Choice struct{
		// Mostly ignore all of this
		Choose int `json:"choose"`
	} `json:"choice"`
	
	// prerequisites
}

func (b Buff)IsNil() bool { return b.Index == "" }

func (b Buff)GetClass() string { return b.Class.Name }
func (b Buff)GetSubclass() string { return b.Subclass.Name }

func (b Buff)String() string { return b.Summary() }

func (b Buff)Type() string { return "buff" }

func (b Buff)Summary() string {
	subclass := ""
	if b.GetSubclass() != "" {
		subclass = "/" + b.GetSubclass()
	}
	return fmt.Sprintf("[%s] %s%s, L%d: %s", b.Index, b.GetClass(), subclass, b.Level, b.Name)
}

func (b Buff)Description() string {
	return fmt.Sprintf("%s\n%s", b.Summary(), b.Desc)
}



type BuffList map[string]Buff

// Implement the lookup interface
func (bl BuffList)Lookup(namelike string) []Entryer {
	if v,exists := bl[namelike]; exists {
		return []Entryer{v}
	}

	ret := []Entryer{}
	for _,v := range bl {
		if strings.Contains(strings.ToLower(v.Name),strings.ToLower(namelike)) {
			ret = append(ret, Entryer(v))
		}
	}
	return ret
}
func (bl BuffList)LookupFirst(namelike string) Entryer {
	if m := bl.Lookup(namelike); len(m) > 0 {
		return m[0]
	}
	return nil
}

// ForClass finds all the buffs that apply for the given character
// class / subclass / level.
func (bl BuffList)ForClass(class, subclass string, level int) []Buff {
	ret := []Buff{}

	for _,buff := range bl {
		if !strings.EqualFold(class, buff.GetClass()) || buff.Level > level {
			continue
		}
		if buff.GetSubclass() != "" && !strings.EqualFold(subclass, buff.GetSubclass()) {
			continue
		}

		// The 'choice' and 'group' elements are things the character picks; don't include here.
		if buff.Group != "" || buff.Choice.Choose > 0 {
			continue
		}

		ret = append(ret, buff)
	}

	return ret
}
