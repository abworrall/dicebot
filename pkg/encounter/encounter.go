package encounter

import(
	"fmt"
	"sort"

	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/roll"
	"github.com/abworrall/dicebot/pkg/rules"
)

type Encounter struct {
	Name string
	Combatants map[string]Combatter
	Init Initiative
	TurnNumber int

	GroupCounts map[string]int   // e.g. GroupCounts["goblin"] = 8 when there are 8 goblins in play
}
func (e Encounter)IsNil() bool { return len(e.Combatants) == 0 }

func NewEncounter() Encounter {
	return Encounter{
		Combatants: map[string]Combatter{},
		Init: NewInitiative(),
		GroupCounts: map[string]int{},
	}
}

func (e Encounter)String() string {
	if e.IsNil() { return "nah, everyone is chillin'" }

	str := ""

	names := []string{}
	for _,v := range e.Combatants {
		names = append(names, v.GetName())
	}
	sort.Strings(names)
	for _,name := range names {
		str += CombatterToString(e.Combatants[name]) + "\n"
	}

	str += "\n" + e.Init.String()

	return str
}

// Returns 1 for monster instance of that name, then 2, etc
func (e *Encounter)NextGroupIndex(name string) int {
	e.GroupCounts[name]++ // auto-inits
	return e.GroupCounts[name]
}

func (e *Encounter)Add(c Combatter) {
	e.Combatants[c.GetName()] = c
}

func (e *Encounter)Lookup(name string) (Combatter, bool) {
	c,exists := e.Combatants[name]
	return c,exists
}

type AttackSpec struct {
	Targets []Combatter

	// Set this a fully auto attack using the attacker's damager
	Attacker Combatter

	// Attacking with a spell, or with a weapon (damager)
	SpellName string
	DamagerName string

	// Various conditions that can apply to an attack
	SpellCastingLevel int // If zero, means the base level of the spell
	WithAdvantage bool
	WithDisadvantage bool
}
func NewAttackSpec() AttackSpec {
	return AttackSpec{Targets: []Combatter{}}
}

func (e *Encounter)Attack(spec AttackSpec) string {
	if (spec.SpellName != "") {
		return e.AttackWithSpell(spec)
	} else {
		return e.AttackWithWeapon(spec)
	}
}

func (e *Encounter)AttackWithWeapon(spec AttackSpec) string {
	c1 := spec.Attacker
	weapon := c1.GetDamager(spec.DamagerName)
	if weapon == nil || weapon.GetName() == "" {
		return fmt.Sprintf("%s did not have a '%s'", c1.GetName(), spec.DamagerName)
	}

	str := ""
	for _,c2 := range spec.Targets {
		str += fmt.Sprintf("%s attacks %s with %s: ", c1.GetName(), c2.GetName(), weapon.GetName())

		// Build the hit roll
		hitRoll := buildHitRoll(spec, weapon, c2)
		hitOutcome := hitRoll.Do()
		str += "Attack - " + hitOutcome.String()

		if hitOutcome.Success {
			// We hit ! Damage roll ?
			if weapon.GetDamageRoll() != "" {
				damageRoll := roll.Parse(weapon.GetDamageRoll())

				if hitOutcome.CriticalHit {
					// Twice as many damage dice !
					str += " DOUBLE DICE"
					damageRoll.NumDice *= 2
				}

				damageOutcome := damageRoll.Do()
				str += " Damage - " + damageOutcome.String()

				c2.AdjustHP(-1 * damageOutcome.Total)

				if hp,_ := c2.GetHP(); hp == 0 {
					str += " they are DEAD"
				}
			}
		}
		str += "\n"
	}

	return str
}


func (e *Encounter)AttackWithSpell(spec AttackSpec) string {
	spell := rules.TheRules.SpellList[spec.SpellName]
	dmg := rules.TheRules.SpellDamageList[spec.SpellName]
	if dmg.IsNil() {
		return fmt.Sprintf("spell '%s' doesn't do attack damage", spec.SpellName)
	}

	str := ""
	
	// We apply spell effects in round-robin form over the targets,
	// until we run out of instances, or it's non-stackable and we run
	// out of targets.
	n,stackable := dmg.Count(spell.Level, spec.SpellCastingLevel)
	for n>0 {
		for _,target := range spec.Targets {
			str += fmt.Sprintf("[%s]", target.GetName())

			if dmg.NeedsAttackRoll() {
				weapon := spec.Attacker.GetDamager("magic")
				hitRoll := buildHitRoll(spec, weapon, target)

				hitOutcome := hitRoll.Do()
				str += " Attack{" + hitOutcome.String() + "}"
				if !hitOutcome.Success {
					str += "\n"
					continue
				}
			}

			damageRoll := dmg.DamageDice(spell.Level, spec.SpellCastingLevel)
			damageOutcome := damageRoll.Do()
			str += fmt.Sprintf(" Damage{%s}", damageOutcome)

			finalDamage := damageOutcome.Total
			if attrStr,effect := dmg.AllowedSave(); attrStr != "" {
				saveRoll := buildSaveRoll(target, attrStr)
				saveOutcome := saveRoll.Do()

				str += fmt.Sprintf(" Save{%s", saveOutcome)

				if saveOutcome.Success {
					finalDamage = int(float64(finalDamage) * effect + 0.5)
					str += fmt.Sprintf(", final damage=%d", finalDamage)
				}
				str += "}"
			}

			target.AdjustHP(-1 * finalDamage)
			if hp,_ := target.GetHP(); hp == 0 {
				str += " they are DEAD"
			}
			str += "\n"

			n--
			if n == 0 {
				// No more instances left to cast
				break
			}
		}
		if !stackable {
			// We've done all targets, and results don't stack, so we're all done
			break
		}
	}

	return str
}

func buildHitRoll(spec AttackSpec, weapon Damager, target Combatter) roll.Roll {
	return roll.Roll{
		NumDice: 1,
		DiceSize: 20,
		Modifier: weapon.GetHitModifier(),
		Target: target.GetArmorClass(),
		WithAdvantage: spec.WithAdvantage,
		WithDisadvantage: spec.WithDisadvantage,
		WithImprovedCritical: spec.Attacker.HasBuff(character.BuffFighterChampionImprovedCritical),
	}
}

func buildSaveRoll(target Combatter, attrStr string) roll.Roll {
	attrKind := character.ParseAttr(attrStr)
	attrVal := target.GetAttr(attrKind)
	attrMod := character.AttrModifier(attrVal)

	saveRoll := roll.Roll{
		NumDice: 1,
		DiceSize: 20,
		Modifier: attrMod,
		Target: 10, // This is the default DC
		Reason: fmt.Sprintf("%s=%d", attrKind, attrVal),
	}

	return saveRoll
}
