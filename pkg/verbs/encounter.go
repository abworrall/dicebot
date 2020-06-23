package verbs

import(
	"fmt"
	"strconv"
	"strings"
	"github.com/abworrall/dicebot/pkg/dnd5e/encounter"
	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
	"github.com/abworrall/dicebot/pkg/roll"
)

// Encounter is stateless as a verb; the encounter data is managed
// by the verb context.
type Encounter struct{}


// Character prep steps:
//   char set weapon longsword  // these populate the list of named attacks you can make
//   char set weapon shortsword
//   char set weapon longsword  // re-establishes the default
//   char set armor scale-mail
//   char set shield 1          // or 0, to disable

// Setup steps:
//   attack -reset
//   attack add goblin.4 wolf.2 bugbear
//   as lanja attack join

// How players make attacks:
//   attack TARGET
//   attack TARGET with WEAPON
//   attack TARGET do 4d6+4
//   attack TARGET1,TARGET2,... do 4d6+4

// How players get attacked:
//   attack TARGET by MONSTER [with ACTION]

func (e Encounter)Help() string { return "[attack join], [attack TARGET [with WEAPON][do DAMAGEROLL]}" }

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
	case "damage":    return e.DoDamage(vc, args[1:])
	}

	// Second, assume it's an attack spec
	attackSpec, err := ParseAttackArgs(vc, args)
	if err != nil {
		return fmt.Sprintf("Problems: %v", err)
	}
	return e.AttackTarget(vc, attackSpec)
}

func (e Encounter)AddPlayer(vc VerbContext) string {
	if vc.Character == nil {
		return "nah"
	}

	c := encounter.NewCombatterFromCharacter(*vc.Character)
	vc.Encounter.Add(c)
	return fmt.Sprintf("%s has joined the fray (now `db roll init`)", vc.Character.Name)
}

func (e Encounter)AddMonsters(vc VerbContext, args []string) string {
	str := ""
	for _,arg := range args {
		str += fmt.Sprintf("adding [%s] ", arg)
		str += e.AddMonster(vc, arg) + "\n"
	}
	return str
}

// goblin
// goblin.3
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
		idx := vc.Encounter.NextGroupIndex(name)
		c := encounter.NewCombatterFromMonster(m, idx)
		vc.Encounter.Add(c)

		if idx == 1 {
			// Roll initiative !
			initVal,initStr := encounter.CombatterRollInitiative(c)
			vc.Encounter.Init.Set(name,initVal)
			str += " " + initStr
		}
	}

	return str
}

func (e Encounter)DoDamage(vc VerbContext, args []string) string {
	str := ""

	if len(args) != 2 { return "damage TARGET {NN or NdM+X}" }

	c,exists := vc.Encounter.Lookup(args[0])
	if !exists {
		return fmt.Sprintf("combatant '%s' not found", args[0])
	}

	damage, err := strconv.Atoi(args[1])

	if err != nil {
		// Probably a dice then
		r := roll.Parse(args[1])
		r.Reason = "damage"
		if r.Err != nil {
			return fmt.Sprintf("I don't get '%s'", args[1])
		}

		o := r.Do()

		str += fmt.Sprintf("%s\n", o)
		damage = o.Total
	}

	c.TakeDamage(damage)
	str += fmt.Sprintf("%s gets %d damage", args[0], damage)

	hp,_ := c.GetHP()
	if hp == 0 {
		str += " and DIES"
	}

	return str
}

type AttackSpec struct {
	Targets []encounter.Combatter
	Attacker encounter.Combatter
	Weapon string
	DamageRoll string
}

// db attack TARGET [with WEAPON] [do DAMAGEROLL] [by MONSTER]
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
