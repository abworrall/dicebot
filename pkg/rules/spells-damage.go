package rules

import(
	"fmt"
	"github.com/abworrall/dicebot/pkg/roll"
)

type SpellDamage struct{
	Index string `json:"index"`

	AttackRoll string `json:"attack_roll"`     // {ranged,melee} for spells that require an attack roll; empty otherwise

	Save struct{
		Attribute string `json:"attribute"`      // con,dex etc; or blank if there is no save
		Effect float64 `json:"effect"`           // Multiplier on the final damage - e.g. 0.5, for half damage
	} `json:"save"`
	
	Multi struct{
		AllInArea bool `json:"all_in_area"`      // Affects all named targets once, e.g. NumTargets = 99
		NumInstances int `json:"num_instances"`  // This damage occurs N times, each hitting the same target or different
		NumTargets int `json:"num_targets"`      // This damage can apply to N different targets
	} `json:"multi"`

	Damage DamageStruct `json:"damage"`

	PerHigherLevel struct {
		ExtraDamageDice string `json:"extra_damage_dice"` // You get an extra damage dice
		ExtraMulti bool `json:"extra_multi"`    // Whatever the multi-instance rules are, you get one more
	} `json:"per_higher_level"`
}

type SpellDamageList map[string]SpellDamage

func (sd SpellDamage)IsNil() bool { return sd.Index == "" }
func (sd SpellDamage)String() string { return sd.Summary() }

func (sd SpellDamage)NeedsAttackRoll() bool { return sd.AttackRoll != "" }

// Count returns how many times the spell does its thing; the bool
// indicates whether the effects stack (i.e. multiple ones can apply
// to the same target). It needs to know both the base level of the
// spell, and the level it is being cast at
func (sd SpellDamage)Count(spellLevel, castingLevel int) (int, bool) {
	if sd.Multi.AllInArea {
		return 999, false
	}

	extra := 0
	if castingLevel > 0 && sd.PerHigherLevel.ExtraMulti == true {
		extra = (castingLevel - spellLevel)
	}

	// If the spell specifies a fixed number of targets (e.g.
	// scorching-ray), then the damage doesn't stack. But if it
	// specifies a fixed number of instances (e.g. magic-missile), then
	// it does stack.
	if sd.Multi.NumTargets > 1 {
		return sd.Multi.NumTargets + extra, false
	} else if sd.Multi.NumInstances > 1 {
		return sd.Multi.NumInstances + extra, true
	}

	return 1, false
}

// DamageRoll returns the roll to be made (per instance). It needs to
// know the base level of the spell, and the level it is being cast at
func (sd SpellDamage)DamageDice(spellLevel, castingLevel int) roll.Roll {
	rBase := roll.Parse(sd.Damage.String())

	if castingLevel > 0 && sd.PerHigherLevel.ExtraDamageDice != "" {
		rExtra := roll.Parse(sd.PerHigherLevel.ExtraDamageDice)
		n := (castingLevel - spellLevel)
		rExtra.NumDice *= n
		rExtra.Modifier *= n

		return rBase.Combine(rExtra)
	}

	return rBase
}

// Save returns how the tareget cna save (or empty if there is no save), and the
// multiplier to apply to the final damage (e.g. 0.5)
func (sd SpellDamage)AllowedSave() (string, float64) {
	return sd.Save.Attribute, sd.Save.Effect
}

func (sd SpellDamage)Summary() string {
	str := fmt.Sprintf("%s: %s", sd.Index, sd.Damage)
	if sd.NeedsAttackRoll() {
		str += ", attack:"+sd.AttackRoll
	}
	if attr,val := sd.AllowedSave(); attr != "" {
		str += fmt.Sprintf(", save:%s x%.1f", attr, val)
	}
	if sd.Multi.NumInstances > 1 {
		str += fmt.Sprintf(", x%d instances", sd.Multi.NumInstances)
	} else if sd.Multi.NumTargets > 1 {
		str += fmt.Sprintf(", x%d targets", sd.Multi.NumTargets)
	}

	if sd.PerHigherLevel.ExtraDamageDice != "" {
		str += fmt.Sprintf(", +%s per level", sd.PerHigherLevel.ExtraDamageDice)
	} else if sd.PerHigherLevel.ExtraMulti == true {
		str += fmt.Sprintf(", +1 multi per level")
	}

	return str
}

