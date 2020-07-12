package character

import(
	"fmt"

	"github.com/abworrall/dicebot/pkg/rules"
	"github.com/abworrall/dicebot/pkg/spells"
)

func (c *Character)IsSpellCaster() bool {
	// Should really be based on class etc
	return c.Slots.Max[1] > 0
}

func (c *Character)MaxSpellLevel() int {
	switch c.Class {
	case "cleric": fallthrough
	case "wizard":
		return (c.Level+1) / 2
	default:
		return 0
	}
}

func (c *Character)MagicString() string {
	if !c.IsSpellCaster() {
		return "You can't do magic :("
	}
	_,desc := c.GetMagicAttackModifier()

	return fmt.Sprintf("\nSpell Attack Modifier: %s\n\n%s\n--- %s", desc, c.GetCastableSpells(), c.Slots)
}

// SpellsAlwaysMemorized lists whatever spells the character has hardwired for whatever reason
func (c *Character)SpellsAlwaysMemorized() spells.Set {
	s := spells.NewSet()

	if c.HasBuff(BuffClericDivineDomainSpells) {
		for _,spell := range rules.TheRules.SpellList.FindMatching(c.Class, c.Subclass, c.MaxSpellLevel()) {
			s.Add(spell.Index)
		}
	}

	return s
}

// GetSpellsMemorized returns a set of all the spells that can be cast right now
func (c *Character)GetCastableSpells() spells.Set {
	return c.SpellsMemorized.UnionWith(c.SpellsAlwaysMemorized(), 2)
}
