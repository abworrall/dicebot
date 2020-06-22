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


// NewCombatterFromMonster takes a monster definition, clones it all,
// and wraps it up so it can act as a Comatter instance in an encounter

func NewCombatterFromMonster(m rules.Monster, i int) Combatter {
	mc := MonsterCombatter{
		CurrentHitPoints: m.HitPoints,
		SpecificName: fmt.Sprintf("%s.%d", m.Index, i),
		Monster: m,
	}

	return mc
}

type MonsterCombatter struct {
	CurrentHitPoints int
	SpecificName     string // e.g. "goblin.1"
	
	rules.Monster
}

func (mc MonsterCombatter)GetName() string { return mc.SpecificName }
func (mc MonsterCombatter)GetGroup() string { return mc.Index }
func (mc MonsterCombatter)GetArmorClass() int { return mc.Monster.ArmorClass }
func (mc MonsterCombatter)GetHP() (int, int) { return mc.CurrentHitPoints, mc.Monster.HitPoints }
func (mc MonsterCombatter)GetDamager(string) Damager { return nil }
func (mc MonsterCombatter)TakeDamage(d int) { mc.CurrentHitPoints -= d }

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
