package verbs

import(
	"fmt"
	"strconv"
	"strings"
	"github.com/abworrall/dicebot/pkg/encounter"
	"github.com/abworrall/dicebot/pkg/rules"
	"github.com/abworrall/dicebot/pkg/roll"
)

// Encounter is stateless as a verb; the encounter data is managed
// by the verb context.
type Encounter struct{}

// Character prep steps:
//   char set weapon longsword   // populate the list of named attacks you can make by adding weapons
//   char set weapon shortsword
//   char set weapon longsword   // default is the last one added (OK to add more than once)
//   char set armor scale-mail
//   char set shield 1           // or 0, to disable

// Setup steps:
//   attack -reset
//   attack add goblin.4 wolf.2 bugbear  // add some friends to the encounter

// How players get involved:
//   attack join
//   attack TARGET
//   attack TARGET TARGET2 TARGET3
//   attack TARGET WEAPON [dis]advantage
//   attack TARGET SPELL
//   attack TARGET SPELL level 5  // which level to cast at
//   attack TARGET 4d6+4

//   attack ... by PLAYER    // one player playing another as an NPC

// On the fly adjustments
//   attack TARGET tweak FIELD VALUE
//   attack TARGET tweak ac +2
//   attack TARGET tweak hp -7

// How players get attacked:
//   attack TARGET by MONSTER [ACTION]

// db attack [by MONSTER] TARGET1 [TARGET2 ...] {WEAPON,SPELL,DAMAGEROLL}] [dis]advantage] [tweak hp 400]

func (e Encounter)Help() string { return "[join], [TARGET ...] {WEAPON,SPELL,ROLL} [advantage] [level NN] [by PLAYER] [tweak {hp,ac} NN]" }

func (e Encounter)Process(vc VerbContext, args []string) string {
	if len(args) < 1 { return vc.Encounter.String() }

	// First, look for named commands
	switch args[0] {
	case "-reset":
		new := encounter.NewEncounter()
		*vc.Encounter = new
		return "encounter reset"

	case "add":       return e.AddMonsters(vc, args[1:])
	case "join":      return e.AddPlayer(vc)
	}

	// Else, assume it's an attack spec
	attack, err := ParseAttackArgs(vc, args)
	if err != nil {
		return fmt.Sprintf("Problems: %v", err)
	}
	return e.AttackTargets(vc, attack)
}

// AddPlayer adds a PC to the encounter. Subsequent calls on the same
// PC reimport the character (i.e. with new weapons), resetting HP,
// but retaining init.
func (e Encounter)AddPlayer(vc VerbContext) string {
	if vc.Character == nil {
		return "nah"
	}

	c := encounter.NewCombatterFromCharacter(*vc.Character)
	vc.Encounter.Add(c)

	initStr := ""
	if initVal := vc.Encounter.Init.Get(vc.Character.Name); initVal > 0 {
		initStr = fmt.Sprintf("%s already has initiative - it is %d", vc.Character.Name, initVal)
	} else {
		initVal,initStr = encounter.CombatterRollInitiative(c)
		vc.Encounter.Init.Set(vc.Character.Name, initVal)
	}

	return fmt.Sprintf("%s has joined the fray (init: %s)", vc.Character.Name, initStr)
}

func (e Encounter)AddMonsters(vc VerbContext, args []string) string {
	str := ""
	for _,arg := range args {
		str += fmt.Sprintf("adding [%s] ", arg)
		str += e.AddMonster(vc, arg) + "\n"
	}
	return str
}

// `nameStr` has an optional count, e.g. `goblin.3`
func (e Encounter)AddMonster(vc VerbContext, nameStr string) string {
	name := nameStr
	bits := strings.Split(nameStr, ".")
	n := 1
	if len(bits) == 2 {
		name = bits[0]
		n,_ = strconv.Atoi(bits[1])
	}

	m,exists := rules.TheRules.MonsterList[name]
	if !exists {
		return fmt.Sprintf("monster '%s' not found", name)
	}

	str := fmt.Sprintf("added %s x%d", name, n)
	
	for i:=0; i<n; i++ {
		// If there are already some monsters of the same type, make sure
		// we starte numbering where they left off.
		idx := vc.Encounter.NextGroupIndex(name)
		c := encounter.NewCombatterFromMonster(m, idx)
		vc.Encounter.Add(c)

		if idx == 1 {
			// We are the first of our kind - roll initiative !
			initVal,initStr := encounter.CombatterRollInitiative(c)
			vc.Encounter.Init.Set(name,initVal)
			str += " " + initStr
		}
	}

	return str
}

type Attack struct {
	encounter.AttackSpec   // Build up the attack spec ...

	TweakHPAmount  int     // ... or tweak a value in the target ...
	TweakACAmount  int
	DamageRoll     string  // ... or jump straight to a damage roll ...
}

// db attack [by MONSTER] TARGET1 [TARGET2 ...] {WEAPON,SPELL,DAMAGEROLL}] [dis]advantage] [level NN] [tweak hp 400]
func ParseAttackArgs(vc VerbContext, args []string) (Attack, error) {
	if len(args) == 0 {
		return Attack{}, fmt.Errorf("not enough args at all")
	}

	attack := Attack{AttackSpec: encounter.NewAttackSpec()}
	attackerName := vc.Character.Name
	if c,exists := vc.Encounter.Lookup(attackerName); ! exists {
		return Attack{}, fmt.Errorf("attacker combatant '%s' not found", attackerName)
	} else {
		attack.AttackSpec.Attacker = c
	}	

	// Keep shifting elements off the front of the slice until it is empty
	for len(args) > 0 {
		if args[0] == "by" {
			if len(args) <2 { return Attack{}, fmt.Errorf("not enough args for 'by PLAYER'") }
			if c,exists := vc.Encounter.Lookup(args[1]); ! exists {
				return Attack{}, fmt.Errorf("attacker combatant '%s' not found", args[1])
				attack.AttackSpec.Attacker = c
			}
			args = args[1:] // Hacky way to keep arg eating in sync since this is a two-arg eat

		} else if args[0] == "level" {
			if len(args) < 2 { return Attack{}, fmt.Errorf("not enough args for level") }
			lvl,_ := strconv.Atoi(args[1])
			attack.AttackSpec.SpellCastingLevel = lvl
			args = args[1:] // Hacky way to keep arg eating in sync since this is a two-arg eat
			
		} else if args[0] == "advantage" {
			attack.AttackSpec.WithAdvantage = true

		} else if args[0] == "disadvantage" {
			attack.AttackSpec.WithDisadvantage = true

		} else if args[0] == "tweak" { // tweak hp -4
			if len(args) < 3 { return Attack{}, fmt.Errorf("not enough args for tweak") }
			mod,_ := strconv.Atoi(args[2])
			switch args[1] {
			case "hp": attack.TweakHPAmount = mod
			case "ac": attack.TweakACAmount = mod
			default: 
				return Attack{}, fmt.Errorf("you can't tweak '%s'", args[1])
			}
			args = args[2:] // Hacky way to keep arg eating in sync since this is a three-arg eat
			
		} else if target,exists := vc.Encounter.Lookup(args[0]); exists {
			attack.AttackSpec.Targets = append(attack.AttackSpec.Targets, target)

		} else if rules.TheRules.IsSpell(args[0]) {
			attack.AttackSpec.SpellName = args[0]

		} else if damager := attack.AttackSpec.Attacker.GetDamager(args[0]); damager != nil && damager.GetName() != "" {
			attack.AttackSpec.DamagerName = args[0]

		} else if dmgRoll := roll.Parse(args[0]); ! dmgRoll.IsNil() {
			attack.DamageRoll = args[0]

		} else {
			return Attack{}, fmt.Errorf("I got lost at '%s'", args[0])
		}

		args = args[1:] // shift off whatever we just processed
	}

	return attack, nil
}

func (e Encounter)AttackTargets(vc VerbContext, attack Attack) string {
	attacker := attack.AttackSpec.Attacker
	if hp,_ := attacker.GetHP(); hp <= 0 {
		return fmt.Sprintf("%s is dead ! such cheating", attacker.GetName())
	}

	// If there is a tweak, just do it
	if attack.TweakHPAmount != 0 {
		str := ""
		for _,target := range attack.AttackSpec.Targets {
			target.AdjustHP(attack.TweakHPAmount)
			str += fmt.Sprintf("%s adjusted HP on %s: %+d", attacker.GetName(), target.GetName(), attack.TweakHPAmount)
			if hp,_ := target.GetHP(); hp <= 0 {
				str += " - target is DEAD"
			}
			str += "\n"
		}
		return str
	}
	if attack.TweakACAmount != 0 {
		str := ""
		for _,target := range attack.AttackSpec.Targets {
			target.AdjustAC(attack.TweakACAmount)
			str += fmt.Sprintf("%s adjusted AC on %s: %+d\n", attacker.GetName(), target.GetName(), attack.TweakACAmount)
		}
		return str
	}

	// If there is an explicit damage roll, just do it
	if attack.DamageRoll != "" {
		str := ""
		for _,target := range attack.AttackSpec.Targets {
			outcome := roll.New(attack.DamageRoll)
			target.AdjustHP(-1 * outcome.Total)
			str += fmt.Sprintf("%s damaged %s: %s", attacker.GetName(), target.GetName(), outcome)
			if hp,_ := target.GetHP(); hp <= 0 {
				str += " - they are DEAD"
			}
			str += "\n"
		}
		return str
	}

	str := vc.Encounter.Attack(attack.AttackSpec) + "\n"

	return str
}
