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
//   attack TARGET with WEAPON
//   attack TARGET do 4d6+4
//   attack TARGET hp -17
//   attack TARGET1,TARGET2,... do 4d6+4

// How players get attacked:
//   attack TARGET by MONSTER [with ACTION]

func (e Encounter)Help() string { return "[join], [TARGET [with WEAPON][do DAMAGEROLL][hp +-NN]" }

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
	attackSpec, err := ParseAttackArgs(vc, args)
	if err != nil {
		return fmt.Sprintf("Problems: %v", err)
	}
	return e.AttackTarget(vc, attackSpec)
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

type AttackSpec struct {
	Targets []encounter.Combatter
	Attacker encounter.Combatter
	Weapon string
	DamageRoll string
	DamageAmount int
}

// db attack TARGET [with WEAPON] [do DAMAGEROLL] [by MONSTER] [hp +-NN]
// TODO: withadvantage, for attack rolls
func ParseAttackArgs(vc VerbContext, args []string) (AttackSpec, error) {
	if len(args) == 0 {
		return AttackSpec{}, fmt.Errorf("not enough args at all")
	}

	spec := AttackSpec{Targets: []encounter.Combatter{}}
	targets := ""
	attacker := vc.Character.Name
	
	targets, args = args[0], args[1:]

	for len(args) >= 2 {
		switch args[0] {
		case "do":   spec.DamageRoll = args[1]
		case "with": spec.Weapon = args[1]
		case "by":   attacker = args[1]
		case "hp":
			amount,_ := strconv.Atoi(args[1])
			spec.DamageAmount = amount
		}
		args = args[2:]
	}

	if c,exists := vc.Encounter.Lookup(attacker); ! exists {
		return AttackSpec{}, fmt.Errorf("attacker combatant '%s' not found", attacker)
	} else {
		spec.Attacker = c
	}	

	for _,target := range strings.Split(targets, ",") {
		if c,exists := vc.Encounter.Lookup(target); ! exists {
			return AttackSpec{}, fmt.Errorf("target combatant '%s' not found", target)
		} else {
			spec.Targets = append(spec.Targets, c)
		}
	}

	return spec, nil
}

func (e Encounter)AttackTarget(vc VerbContext, spec AttackSpec) string {
	if hp,_ := spec.Attacker.GetHP(); hp <= 0 {
		return fmt.Sprintf("%s is dead ! such cheating", spec.Attacker.GetName())
	}

	// If there is explicit damage, just apply it
	if spec.DamageAmount != 0 {
		str := ""
		for _,target := range spec.Targets {
			target.TakeDamage(spec.DamageAmount)
			str += fmt.Sprintf("%s damaged %s: %+d", spec.Attacker.GetName(), target.GetName(), spec.DamageAmount)
			if hp,_ := target.GetHP(); hp <= 0 {
				str += " - they are DEAD"
			}
			str += "\n"
		}
		return str
	}

	// If there is an explicit damage roll, just do it
	if spec.DamageRoll != "" {
		str := ""
		for _,target := range spec.Targets {
			outcome := roll.New(spec.DamageRoll)
			target.TakeDamage(outcome.Total)
			str += fmt.Sprintf("%s damaged %s: %s", spec.Attacker.GetName(), target.GetName(), outcome)
			if hp,_ := target.GetHP(); hp <= 0 {
				str += " - they are DEAD"
			}
			str += "\n"
		}
		return str
	}

	if len(spec.Targets) != 1 {
		return "You can only attack one target with your weapon"
	}
	str := vc.Encounter.Attack(spec.Attacker, spec.Targets[0], spec.Weapon) + "\n"

	return str
}
