package encounter

import(
	"encoding/gob"
	"fmt"

	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)

func init() {
	gob.Register(CharacterCombatter{})
}

type CharacterCombatter struct {
	Counts           map[string]int // We abstract any values we need to mutate; we can't pass mc by pointer
	character.Character
}

// NewCombatterFromMonster takes a monster definition, clones it all,
// and wraps it up so it can act as a Comatter instance in an encounter
func NewCombatterFromCharacter(c character.Character) Combatter {
	cc := CharacterCombatter{
		Counts: map[string]int{"hp": c.CurrHitpoints},
		Character: c,
	}

	return cc
}

func (cc CharacterCombatter)GetName() string { return cc.Name }
func (cc CharacterCombatter)GetGroup() string { return cc.Name }
func (cc CharacterCombatter)GetHP() (int, int) { return cc.Counts["hp"], cc.Character.MaxHitpoints }

func (cc CharacterCombatter)TakeDamage(d int) {
	cc.Counts["hp"] -= d
	if cc.Counts["hp"] < 0 {
		cc.Counts["hp"] = 0
	}
}

func (cc CharacterCombatter)GetArmorClass() int {
	ac := 10

	dexBonus := cc.GetModifier(character.Dex)
	
	if cc.Character.Armor != "" {
		armor := rules.TheRules.EquipmentList[cc.Character.Armor]
		ac = armor.ArmorClass.Base
		if !armor.ArmorClass.DexBonus {
			// Not allowed with this kind of armor
			dexBonus = 0
		}
	}

	ac += dexBonus
	if cc.Character.Shield {
		ac += 2
	}

	return ac
}

func (cc CharacterCombatter)GetDamagerNames() []string {
	ret := []string{}
	for k,_ := range cc.Character.Weapons {
		ret = append(ret, k)
	}
	return ret
}

func (cc CharacterCombatter)GetDamager(name string) Damager {
	if name == "" {
		name = cc.Character.CurrWeapon
	}
	if name == "" {
		return nil
	}
	if _,exists := cc.Character.Weapons[name]; !exists {
		return nil
	}
	
	w := rules.TheRules.EquipmentList[name]
	if w.IsNil() {
		return nil
	}

	return WeaponDamager{
		Item: w,
		Character: cc.Character,
	}
}


func (cc CharacterCombatter)GetAttr(k character.AttrKind) int {
	switch k {
	case character.Str: return cc.Str
	case character.Int: return cc.Int
	case character.Wis: return cc.Wis
	case character.Con: return cc.Con
	case character.Cha: return cc.Cha
	case character.Dex: return cc.Dex
	case character.Per: return cc.Per
	default:  return -1
	}
}


// WeaponDamager wraps up a weapon object (from rules) as a Damager
type WeaponDamager struct {
	character.Character
	rules.Item
}

func (wd WeaponDamager)	GetName() string {
	return wd.Item.Index
}

func (wd WeaponDamager)	GetHitModifier() int {
	strMod := wd.Character.GetModifier(character.Str)
	dexMod := wd.Character.GetModifier(character.Dex)

	mod := 0

	// First, basic attr bonus (melee or ranged)
	switch wd.Item.WeaponRange {
	case "Melee":  mod = strMod
	case "Ranged": mod = dexMod
	}

	if wd.Item.HasProperty("Finesse") {
		if strMod > dexMod {
			mod = strMod
		} else {
			mod = dexMod
		}
	}
	
	return mod
}

func (wd WeaponDamager)	GetDamageRoll() string {
	str := wd.Item.Damage.DamageDice

	bonus := wd.Item.Damage.DamageBonus  // maybe magic items have this ?
	bonus += wd.GetHitModifier()         // what goes for attack, goes for damage
	
	if bonus != 0 {
		str += fmt.Sprintf("%+d", bonus)
	}

	return str
}

