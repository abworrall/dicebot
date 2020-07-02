package character

import "fmt"

// Various class-specific guff - subtypes, class features, etc.

type Buff string

const(
	// If you add a new one, also add to the switch in AddBuff below
	BuffFighterChampionImprovedCritical  Buff = "improved-critical"
	BuffFighterFightingStyleDefense      Buff = "defense"  // Adds +1 to AC
)

func (c *Character)HasBuff(b Buff) bool {
	autobuffs := c.AutoBuffs()

	_,exists1 := c.Buffs[b]
	_,exists2 := autobuffs[b]

	return exists1 || exists2
}

func (c *Character)AddBuff(b Buff) error {
	switch b {
	case BuffFighterChampionImprovedCritical: fallthrough
	case BuffFighterFightingStyleDefense:
		c.Buffs[b] = 1
		return nil

	default:
		return fmt.Errorf("Buff '%s' not known", b)
	}
}

// Autobuffs returns a list of the buffs you get automatically just
// because of class/subclass/level/race.
func (c *Character)AutoBuffs() map[Buff]int {
	m := map[Buff]int{}

	// Consider rebasing on top of https://github.com/bagelbits/5e-database/blob/master/src/5e-SRD-Features.json
	if c.Class == "fighter" {
		if c.Subclass == "champion" {
			if c.Level >= 3 { m[BuffFighterChampionImprovedCritical] = 1 }
		}
	}

	return m
}
