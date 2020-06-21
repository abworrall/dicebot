package rules

import(
	"fmt"
	"strings"
)

type Item struct {
	Index string `json:"index"`
	Name string `json:"name"`
	EquipmentCategory struct {
		Name string `json:"name"`
	} `json:"equipment_category"`

	// There is also vehicle_category, etc; only parsing weapony fields for now.
	WeaponCategory string `json:"weapon_category"`
	WeaponRange string `json:"weapon_range"`
	CategoryRange string `json:"category_range"`

	Damage DamageStruct `json:"damage"`

	Range struct {
		Normal float64 `json:"normal"`
		Long float64 `json:"long"`
	} `json:"range"`

	Weight float64 `json:"weight"`
	
	Properties []struct{
		Name string `json:"name"`
	} `json:"properties"`

	Damage2H DamageStruct `json:"2h_damage"`
}

// Use named struct, since it occurs twice in the data (where two-handed damage differs)
type DamageStruct struct {
	DamageDice string `json:"damage_dice"`
	DamageBonus float64 `json:"damage_bonus"`
	DamageType struct {
		Name string `json:"name"`
	} `json:"damage_type"`
}

func (d DamageStruct)String() string {
	s := d.DamageDice
	if d.DamageBonus != 0 {
		s += fmt.Sprintf("%+d", int(d.DamageBonus))
	}
	return s
}

func (i Item)String() string { return i.Summary() }

func (i Item)Type() string { return "equipment" }

func (i Item)Summary() string {
	s := fmt.Sprintf("%s: ", i.Index)

	if i.WeaponCategory != "" {
		// It's a weapon !
		s += "damage:"+i.Damage.String()

		if i.WeaponRange == "Ranged" {
			s += fmt.Sprintf(" range[%.0f,%.0f]", i.Range.Normal, i.Range.Long)
		}
		
		if len(i.Properties) > 0 {
			props := make([]string, len(i.Properties))
			for j,p := range i.Properties {
				props[j] = p.Name
			}
			s += fmt.Sprintf(" [%s]", strings.Join(props, ","))
		}
	}
	return s
}

func (i Item)Description() string { return i.Summary() + "\n" }

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
