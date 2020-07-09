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
	_,autoStr := c.SpellsAlwaysMemorized()

	return fmt.Sprintf("\nSpell Attack Modifier: %s\n\n%s%s\n--- %s", desc, c.SpellsMemorized, autoStr, c.Slots)
}

// SpellsAlwaysMemorized lists whatever spells the character can always cast
func (c *Character)SpellsAlwaysMemorized() (spells.Set, string) {
	s := spells.NewSet()
	str := ""
	
	for _,spell := range rules.TheRules.SpellList.FindMatching(c.Class, c.Subclass, c.MaxSpellLevel()) {
		s.Add(spell.Index)
		str += fmt.Sprintf(" L%d {%s} cast:{%s}\n", spell.Level, spell.Index, spell.CastingTime)
	}
	
	return s, str
}
