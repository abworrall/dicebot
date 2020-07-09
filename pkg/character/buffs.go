package character

import(
	"fmt"
	"github.com/abworrall/dicebot/pkg/rules"
)

// Various class-specific guff - subtypes, class features, etc.

const(
	// This list is for buffs that we make use of elsewhere in the bot 
	BuffFighterChampionImprovedCritical = "improved-critical"
	BuffFighterFightingStyleDefense = "fighter-fighting-style-defense"
	BuffFighterFightingStyleDueling = "fighter-fighting-style-dueling"
)

func (c *Character)AddBuff(b string) error {
	if _,exists := rules.TheRules.BuffList[b]; !exists {
		return fmt.Errorf("Buff '%s' not known", b)
	}

	c.Buffs[b] = 1
	return nil
}


func (c *Character)RemoveBuff(b string) error {
	if _,exists := c.Buffs[b]; !exists {
		return fmt.Errorf("You don't have buff '%s'", b)
	}

	delete(c.Buffs, b)

	return nil
}

func (c *Character)HasBuff(b string) bool {
	autobuffs := c.AutoBuffs()

	_,exists1 := c.Buffs[b]
	_,exists2 := autobuffs[b]

	return exists1 || exists2
}

// Autobuffs returns a list of the buffs you get automatically just
// because of class/subclass/level/race.
func (c *Character)AutoBuffs() map[string]int {
	m := map[string]int{}

	ruleBuffs := rules.TheRules.BuffList.ForClass(c.Class, c.Subclass, c.Level)
	for _,ruleBuff := range ruleBuffs {
		m[ruleBuff.Index] = 1
	}

	return m
}
