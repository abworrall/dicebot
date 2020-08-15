package character

import(
	"fmt"
	"strings"

	"github.com/abworrall/dicebot/pkg/rules"
)

// Lookup tables and logic for figuring out various level/class/proficiency based modifiers

func XPAtLevel(lvl int) int {
	vals := []int{
		0, 300, 900, 2700, 6500, 14000, 23000, 34000, 48000, 64000,
		85000, 100000, 120000, 140000, 165000, 195000, 225000, 265000, 305000, 355000,
	}
	if lvl < 1 || lvl > 20 { return -1 }
	return vals[lvl-1]
}

func (c *Character)NextLevelAt() int {
	return XPAtLevel(c.Level + 1)
}

func AttrModifier(attrVal int) int {
	if attrVal < 1 { attrVal = 1 }
	if attrVal > 30 { attrVal = 30 }

	// This table is pretty simple.
	return attrVal/2 - 5
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
		frags = append(frags, fmt.Sprintf("{%s; base AC=%d}", armor.Index, ac))

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

		// This only applies when wearing armor
		if c.HasBuff(BuffFightingStyleDefense) {
			ac += 1
			frags = append(frags, fmt.Sprintf("{%s %+d}", BuffFightingStyleDefense, 1))
		}
		
	} else {
		frags = append(frags, "{no armor; base AC=10}")
		if dexBonus > 0 {
			frags = append(frags, fmt.Sprintf("{dex bonus %+d}", dexBonus))
		}
	}

	ac += dexBonus
	if c.Shield {
		ac += 2
		frags = append(frags, fmt.Sprintf("{shield %+d}", 2))
	}

	frags = append(frags, fmt.Sprintf("{final AC=%d}", ac))

	return ac, strings.Join(frags, " ")
}

func (c *Character)GetWeaponDamageRoll(w rules.Item) string {
	str := w.Damage.DamageDice

	bonus := w.Damage.DamageBonus  // maybe magic items have this ?
	mod,_ := c.GetWeaponDamageModifier(w)
	bonus += mod

	if bonus != 0 {
		str += fmt.Sprintf("%+d", bonus)
	}

	return str
}

// https://rpg.stackexchange.com/questions/72910/how-do-i-figure-the-dice-and-bonuses-for-attack-rolls-and-damage-rolls

// For when the character casts a spell, and needs a 'magic attack roll'
func (c *Character)GetMagicAttackModifier() (int, string) {
	attr := c.GetSpellcastingAbilityAttr()
	mod := c.GetModifier(attr)
	frags := []string{fmt.Sprintf("{spellcasting attr %s %+d}", attr, mod)}

	// Then, proficiency bonus !
	mod += c.ProficiencyBonus()
	frags = append(frags, fmt.Sprintf("{proficiency %+d}", c.ProficiencyBonus()))
	frags = append(frags, fmt.Sprintf("{total %+d}", mod))

	return mod, strings.Join(frags, " ")
}

func (c *Character)GetWeaponAttackModifier(w rules.Item) (int, string) {
	// Start off with the basic ability modifier
	mod, desc := c.GetWeaponAbilityModifier(w)
	frags := []string{desc}

	// Then, proficiency bonus !
	mod += c.ProficiencyBonus()
	frags = append(frags, fmt.Sprintf("{proficiency %+d}", c.ProficiencyBonus()))
	frags = append(frags, fmt.Sprintf("{total=%+d}", mod))

	return mod, strings.Join(frags, " ")
}

func (c *Character)GetWeaponDamageModifier(w rules.Item) (int, string) {
	// Start off with the basic ability modifier
	mod, desc := c.GetWeaponAbilityModifier(w)
	frags := []string{desc}

	// The Dueling buffs give a +2 damage bonus
	if w.WeaponRange == "Melee" && c.HasBuff(BuffFightingStyleDueling) {
		frags = append(frags, fmt.Sprintf("{%s %+d}", BuffFightingStyleDueling, 2))
		mod += 2
	}
	
	// No proficiency bonus for damage.

	frags = append(frags, fmt.Sprintf("{total=%+d}", mod))
	
	return mod, strings.Join(frags, " ")
}

// GetWeaponAbilityModifier computes the modifier deriving from the character's
// ability for the given weapon.
func (c *Character)GetWeaponAbilityModifier(w rules.Item) (int, string) {
	strMod := c.GetModifier(Str)
	dexMod := c.GetModifier(Dex)

	frags := []string{}
	mod := 0

	if w.HasProperty("Finesse") {
		if strMod > dexMod {
			frags = append(frags, fmt.Sprintf("{finesse; str %+d}", strMod))
			mod = strMod
		} else {
			frags = append(frags, fmt.Sprintf("{finesse; dex %+d}", dexMod))
			mod = dexMod
		}

	} else {
		// First, basic attr bonus (melee or ranged)
		switch w.WeaponRange {
		case "Melee":
			mod = strMod
			frags = append(frags, fmt.Sprintf("{melee; str %+d}", strMod))

		case "Ranged":
			mod = dexMod
			frags = append(frags, fmt.Sprintf("{ranged; dex %+d}", dexMod))
		}
	}

	return mod, strings.Join(frags, " ")
}
