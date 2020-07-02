package encounter

import(
	"encoding/gob"

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
	ac,_ := cc.Character.GetArmorClass()
	return ac
}

func (cc CharacterCombatter)GetDamagerNames() []string {
	ret := []string{}
	for k,_ := range cc.Character.Weapons {
		ret = append(ret, k)
	}

	if cc.Character.IsSpellCaster() {
		ret = append(ret, "magic")
	}

	return ret
}

func (cc CharacterCombatter)GetDamager(name string) Damager {
	if name == "magic" {
		return MagicDamager{Character: cc.Character}
	}

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

func (cc CharacterCombatter)HasBuff(b string) bool { return cc.Character.HasBuff(b) }

// WeaponDamager wraps up a weapon object (from rules) as a Damager
type WeaponDamager struct {
	character.Character
	rules.Item
}
func (wd WeaponDamager)	GetName() string { return wd.Item.Index }
func (wd WeaponDamager)	GetDamageRoll() string { return wd.Character.GetWeaponDamageRoll(wd.Item) }
func (wd WeaponDamager)	GetHitModifier() int {
	mod,_ := wd.Character.GetWeaponAttackModifier(wd.Item)
	return mod
}

// MagicDamager is a shim to represent a magic attack; it has no
// damage, since that's all handed separately.
type MagicDamager struct {
	character.Character
}
func (md MagicDamager)GetName() string { return "magic" }
func (md MagicDamager)GetDamageRoll() string { return "" }
func (md MagicDamager)GetHitModifier() int {
	mod,_ := md.Character.GetMagicAttackModifier()
	return mod
}


