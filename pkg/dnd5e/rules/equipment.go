package rules

import(
	"fmt"
	"strings"
)

type Item struct {
	Index string `json:"index"`
	Name string `json:"name"`

	Cost struct {
		Quantity int `json:"quantity"`
		Unit string `json:"unit"`
	} `json:"cost"`
	
	EquipmentCategory struct {
		Name string `json:"name"`
	} `json:"equipment_category"`

	// There is also vehicle_category, etc; only parsing weapony fields for now.
	WeaponCategory string `json:"weapon_category"`
	WeaponRange string `json:"weapon_range"`
	CategoryRange string `json:"category_range"`

	Damage DamageStruct `json:"damage"`

	Range struct {
		Normal int `json:"normal"`
		Long int `json:"long"`
	} `json:"range"`

	ArmorCategory string `json:"armor_category"`
	ArmorClass struct {
		Base int `json:"base"`
		DexBonus bool `json:"dex_bonus"`
		MaxBonus int `json:"max_bonus"`
	} `json:"armor_class"`

	Weight int `json:"weight"`
	Desc []string `json:"desc"`
	
	Properties []struct{
		Name string `json:"name"`
	} `json:"properties"`

	Damage2H DamageStruct `json:"2h_damage"`
}

// Use named struct, since it occurs twice in the data (where two-handed damage differs)
type DamageStruct struct {
	DamageDice string `json:"damage_dice"`
	DamageBonus int `json:"damage_bonus"`
	DamageType struct {
		Name string `json:"name"`
	} `json:"damage_type"`
}

func (i Item)IsNil() bool { return i.Index == "" }

func (d DamageStruct)String() string {
	s := d.DamageDice
	if d.DamageBonus != 0 {
		s += fmt.Sprintf("%+d", d.DamageBonus)
	}
	return s
}

func (i Item)String() string { return i.Summary()}

func (i Item)Type() string { return "equipment" }

func (i Item)Summary() string {
	s := fmt.Sprintf("%s: ", i.Index)

	if i.WeaponCategory != "" {
		// It's a weapon !
		s += i.WeaponDamageString()

		if i.WeaponRange == "Ranged" {
			s += fmt.Sprintf(" range[%d,%d]", i.Range.Normal, i.Range.Long)
		}

		if len(i.Properties) > 0 {
			props := make([]string, len(i.Properties))
			for j,p := range i.Properties {
				props[j] = p.Name
			}
			s += fmt.Sprintf(" [%s]", strings.Join(props, ","))
		}

	} else if i.ArmorCategory != "" {
		// It's a kind of armor !
		s += fmt.Sprintf("AC:%d, dex-bonus:%v, max-bonus:%+d",
			i.ArmorClass.Base, i.ArmorClass.DexBonus, i.ArmorClass.MaxBonus)
	}

	s += fmt.Sprintf(" {%d%s, %dlb}", i.Cost.Quantity, i.Cost.Unit, i.Weight)
	
	return s
}

func (i Item)WeaponDamageString() string {
	if i.WeaponCategory == "" { return "" }
	s := "damage:"+i.Damage.String()
	if i.Damage2H.DamageDice != "" {
		s += " (2h_damage:"+i.Damage2H.String()+")"
	}
	return s
}

func (i Item)Description() string { return i.Summary() + "\n" + i.GetDescriptions() }

func (i Item)HasProperty(name string) bool {
	for _,prop := range i.Properties {
		if strings.ToLower(name) == strings.ToLower(prop.Name) {
			return true
		}
	}
	return false
}

func (i Item)GetDescriptions() string {
	str := ""
	for _,d := range i.Desc {
		str += "[" + d + "]\n"
	}
	return str
}

// EquipmentList simply maps the `Index` of each item to the full object
type EquipmentList map[string]Item

// Implement the lookup interface
func (el EquipmentList)Lookup(namelike string) []Entryer {
	if v,exists := el[namelike]; exists {
		return []Entryer{v}
	}

	ret := []Entryer{}
	for _,v := range el {
		if strings.Contains(strings.ToLower(v.Name),strings.ToLower(namelike)) {
			ret = append(ret, Entryer(v))
		}
	}
	return ret
}

func (el EquipmentList)LookupFirst(namelike string) Entryer {
	if m := el.Lookup(namelike); len(m) > 0 {
		return m[0]
	}
	return nil
}
