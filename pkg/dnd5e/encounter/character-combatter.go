package encounter

import(
	"encoding/gob"
	"github.com/abworrall/dicebot/pkg/character"
)

func init() {
	gob.Register(CharacterCombatter{})
}

type CharacterCombatter struct {
	CurrentHitPoints int
	character.Character
}

// NewCombatterFromMonster takes a monster definition, clones it all,
// and wraps it up so it can act as a Comatter instance in an encounter
func NewCombatterFromCharacter(c character.Character) Combatter {
	cc := CharacterCombatter{
		CurrentHitPoints: c.CurrHitpoints,
		Character: c,
	}

	return cc
}

func (cc CharacterCombatter)GetName() string { return cc.Name }
func (cc CharacterCombatter)GetGroup() string { return cc.Name }
func (cc CharacterCombatter)GetHP() (int, int) { return cc.CurrentHitPoints, cc.Character.MaxHitpoints }
func (cc CharacterCombatter)GetDamager(string) Damager { return nil }
func (cc CharacterCombatter)TakeDamage(d int) { cc.CurrentHitPoints -= d }

// TODO: armor ?!
func (cc CharacterCombatter)GetArmorClass() int { return 10 }


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
