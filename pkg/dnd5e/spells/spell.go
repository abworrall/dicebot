package spells

import(
	"fmt"
	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)

// We don't represent spells as such; the rules objects have all that

// LookupSpell pulls a spell out of the rules (maybe test result with `isNil`)
func Lookup(name string) rules.Spell {
	return rules.TheRules.SpellList[name]
}

// Cast verifies that the spell is in the Set can be cast at the given
// `level`, and if so, consumes the Slot and returns `nil`.
func Cast(set Set, slots *Slots, name string, level int) error {
	spell := Lookup(name)

	if spell.IsNil() {
		return fmt.Errorf("spell '%s' is not even a known spell", name)
	} else if !set.Contains(name) {
		return fmt.Errorf("spell '%s' is not in the current spellset", name)
	}

	// Treat level 0 as default spell level
	if level == 0 {
		level = spell.Level
	} else if level < spell.Level {
		return fmt.Errorf("spell '%s' is L%d, can't use an L%d slot", spell.Level, level)
	}

	if slots.Curr[level] <= 0 {
		return fmt.Errorf("you don't have any L%d slots left", level)
	}

	slots.Curr[level]--

	return nil
}
