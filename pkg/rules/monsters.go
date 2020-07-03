package rules

import(
	"fmt"
	"strings"
)

type Monster struct {
	Index string `json:"index"`
	Name string `json:"name"`
	Size string `json:"size"`
	Alignment string `json:"alignment"`

	ArmorClass int `json:"armor_class"`
	HitPoints int `json:"hit_points"`
	HitDice string `json:"hit_dice"`

	Speed struct {
		Walk string `json:"walk"`
		Swim string `json:"swim"`
		Fly string `json:"fly"`
	} `json:"speed"`

	Str int `json:"strength"`
	Dex int `json:"dexterity"`
	Con int `json:"constitution"`
	Int int `json:"intelligence"`
	Wis int `json:"wisdom"`
	Cha int `json:"charisma"`
	
	Abilities []struct{
		Name string `json:"name"`
		Desc string `json:"desc"`
	} `json:"special_abilities"`

	Actions []ActionStruct `json:"actions"`
}

type ActionStruct struct{
	Index string
	Name string `json:"name"`
	Desc string `json:"desc"`
	AttackBonus int `json:"attack_bonus"`
	Damage []DamageStruct `json:"damage"`
}

func (m Monster)String() string { return m.Summary() }

func (m Monster)Type() string { return "monster" }

func (m Monster)Summary() string {
	return fmt.Sprintf("[%s] AC:%d, HP:%d(%s)", m.Index, m.ArmorClass, m.HitPoints, m.HitDice)
}

func (m Monster)Description() string {
	s := m.Summary() + "\n--Abilities--\n"

	for _,a := range m.Abilities {
		s += fmt.Sprintf("%s: [%s]\n", a.Name, a.Desc)
	}
	s += "--Actions--\n"

	for _,a := range m.Actions {
		damages := make([]string, len(a.Damage))
		for i,d := range a.Damage {
			damages[i] = d.String()
		}
		s += fmt.Sprintf("%s (attack:%+d; damages:[%s]): [%s]\n",
			a.Name, a.AttackBonus, strings.Join(damages, ","), a.Desc)
	}

	return s
}

func (m *Monster)PostLoadFixups() {
	// The actions don't have an index field, but we need to refer to them from outside, so ...
	for i,_ := range m.Actions {
		idx := strings.ToLower(m.Actions[i].Name)
		idx = strings.ReplaceAll(idx, " ", "-")
		m.Actions[i].Index = idx
	}
}


type MonsterList map[string]Monster

// Implement the lookup interface
func (ml MonsterList)Lookup(namelike string) []Entryer {
	if v,exists := ml[namelike]; exists {
		return []Entryer{v}
	}

	ret := []Entryer{}
	for _,v := range ml {
		if strings.Contains(strings.ToLower(v.Name),strings.ToLower(namelike)) {
			ret = append(ret, Entryer(v))
		}
	}
	return ret
}

func (ml MonsterList)LookupFirst(namelike string) Entryer {
	if m := ml.Lookup(namelike); len(m) > 0 {
		return m[0]
	}
	return nil
}
