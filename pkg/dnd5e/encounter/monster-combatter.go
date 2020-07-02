package encounter

import(
	"encoding/gob"
	"fmt"
	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)


func init() {
	gob.Register(MonsterCombatter{})
}

type MonsterCombatter struct {
	Counts           map[string]int // We abstract any values we need to mutate; we can't pass mc by pointer
	SpecificName     string // e.g. "goblin.1"
	
	rules.Monster
}

// NewCombatterFromMonster takes a monster definition, clones it all,
// and wraps it up so it can act as a Comatter instance in an encounter

func NewCombatterFromMonster(m rules.Monster, i int) Combatter {
	mc := MonsterCombatter{
		Counts: map[string]int{"hp": m.HitPoints},
		SpecificName: fmt.Sprintf("%s.%d", m.Index, i),
		Monster: m,
	}

	return mc
}

func (mc MonsterCombatter)TakeDamage(d int) {
	mc.Counts["hp"] -= d
	if mc.Counts["hp"] < 0 {
		mc.Counts["hp"] = 0
	}

}

func (mc MonsterCombatter)GetName() string { return mc.SpecificName }
func (mc MonsterCombatter)GetGroup() string { return mc.Index }
func (mc MonsterCombatter)GetArmorClass() int { return mc.Monster.ArmorClass }
func (mc MonsterCombatter)GetHP() (int, int) { return mc.Counts["hp"], mc.Monster.HitPoints }

func (mc MonsterCombatter)GetAttr(k character.AttrKind) int {
	switch k {
	case character.Str: return mc.Str
	case character.Int: return mc.Int
	case character.Wis: return mc.Wis
	case character.Con: return mc.Con
	case character.Cha: return mc.Cha
	case character.Dex: return mc.Dex
	default:  return -1
	}
}

func (mc MonsterCombatter)HasBuff(b string) bool { return false }

func (mc MonsterCombatter)GetDamagerNames() []string {
	ret := make ([]string, len(mc.Monster.Actions))
	for i,_ := range mc.Monster.Actions {
		ret[i] = mc.Monster.Actions[i].Index
	}
	return ret
}

func (mc MonsterCombatter)GetDamager(name string) Damager {
	// If none specified, but there is only one action anyway, use it
	if name == "" && len(mc.Monster.Actions) == 1 {
		return ActionDamager{Action:mc.Monster.Actions[0], DamageIndex:0}
	}

	for _,action := range mc.Monster.Actions {
		if action.Index == name {
			return ActionDamager{Action:action, DamageIndex:0}
		}
	}
	return nil
}

// ActionDamager wraps up a monster action (from rules) as a Damager
type ActionDamager struct {
	DamageIndex int // TODO: Handle multiple damage entries under one action ?
	Action rules.ActionStruct
}

func (ad ActionDamager)	GetName() string {
	return ad.Action.Index
}

func (ad ActionDamager)	GetHitModifier() int {
	return ad.Action.AttackBonus
}

func (ad ActionDamager)	GetDamageRoll() string {
	str := ad.Action.Damage[ad.DamageIndex].DamageDice
	if ad.Action.Damage[ad.DamageIndex].DamageBonus != 0 {
		str += fmt.Sprintf("%+d", ad.Action.Damage[ad.DamageIndex].DamageBonus)
	}
	return str
}
