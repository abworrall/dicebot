package character

import(
	"fmt"
	"strings"

	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)

// Lookup tables and logic for figuring out various level/class/proficiency based modifiers

func AttrModifier(attrVal int) int {
	// Kinda horrible lookup table
	var attrModifiers = []int{
		0,-5,-4,-4,          // attr scores 0-3
		-3,-3,-2,-2,-1,-1,   // attr scores 4-9
		0,0,1,1,2,2,         // attr scores 10-15
		3,3,4,4,5,5,         // attr scores 16-21
		6,6,7,7,8,8,9,9,10,  // attr scores 22-30
	}
	if attrVal<0 || attrVal>len(attrModifiers) {
		return 0
	}
	return attrModifiers[attrVal]
}

func (c *Character)ProficiencyBonus() int {
	lvl := c.Level
	if lvl < 1 { lvl = 1 }
	if lvl > 20 { lvl = 20 }

	// This table is pretty simple.
	return (lvl-1)/4 + 2
}

// Compute a character's AC. The string returns a description.
func (c *Character)GetArmorClass() (int, string) {
	frags := []string{}

	ac := 10
	dexBonus := c.GetModifier(Dex)
	
	if c.Armor != "" {
		armor := rules.TheRules.EquipmentList[c.Armor]
		ac = armor.ArmorClass.Base
		frags = append(frags, fmt.Sprintf("{%s; AC=%d}", armor.Index, ac))

		// Apply various restrictions on the dex bonus
		if !armor.ArmorClass.DexBonus {
			dexBonus = 0
			frags = append(frags, fmt.Sprintf("{no dex bonus allowed}"))
		} else {
			if armor.ArmorClass.MaxBonus > 0 && dexBonus > armor.ArmorClass.MaxBonus {
				dexBonus = armor.ArmorClass.MaxBonus
				frags = append(frags, fmt.Sprintf("{dex bonus capped at %+d}", dexBonus))
			} else {
				frags = append(frags, fmt.Sprintf("{dex bonus %+d}", dexBonus))
			}
		}

	} else {
		frags = append(frags, "{no armor; AC=10}")
	}

	ac += dexBonus
	if c.Shield {
		ac += 2
		frags = append(frags, fmt.Sprintf("{shield %+d}", 2))
	}

	frags = append(frags, fmt.Sprintf("{final AC:%d}", ac))

	return ac, strings.Join(frags, " ")
}

// For when the character casts a spell, and needs a 'magic attack roll'
func (c *Character)GetMagicAttackModifier() (int, string) {
	attr := c.GetSpellcastingAbilityAttr()
	mod := c.GetModifier(attr)
	frags := []string{fmt.Sprintf("{spellcasting attr %s:%+d}", attr, mod)}

	// Then, proficiency bonus !
	mod += c.ProficiencyBonus()
	frags = append(frags, fmt.Sprintf("{proficiency: %+d}", c.ProficiencyBonus()))


	frags = append(frags, fmt.Sprintf("{total:%+d}", mod))

	return mod, strings.Join(frags, " ")
}

func (c *Character)GetWeaponAttackModifier(w rules.Item) (int, string) {
	strMod := c.GetModifier(Str)
	dexMod := c.GetModifier(Dex)

	frags := []string{}
	mod := 0

	if w.HasProperty("Finesse") {
		if strMod > dexMod {
			frags = append(frags, fmt.Sprintf("{finesse; picked str:%+d}", strMod))
			mod = strMod
		} else {
			frags = append(frags, fmt.Sprintf("{finesse; picked dex:%+d}", strMod))
			mod = dexMod
		}

	} else {
		// First, basic attr bonus (melee or ranged)
		switch w.WeaponRange {
		case "Melee":
			mod = strMod
			frags = append(frags, fmt.Sprintf("{melee; str:%+d}", strMod))

		case "Ranged":
			mod = dexMod
			frags = append(frags, fmt.Sprintf("{ranged; dex:%+d}", dexMod))
		}
	}

	// Then, proficiency bonus !
	mod += c.ProficiencyBonus()
	frags = append(frags, fmt.Sprintf("{proficiency: %+d}", c.ProficiencyBonus()))
	
	return mod, strings.Join(frags, " ")
}
