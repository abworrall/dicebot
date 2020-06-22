package verbs

import(
	"fmt"
	"strconv"
	"strings"
	"github.com/abworrall/dicebot/pkg/dnd5e/encounter"
	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
)

// Encounter is stateless as a verb; the encounter data is managed
// by the verb context.
type Encounter struct{}

// attack setup goblin.4 pc:lanja
// attack add goblin.2

func (e Encounter)Help() string { return "oh lordy - use the source, luke" }

func (e Encounter)Process(vc VerbContext, args []string) string {
	if len(args) < 1 { return vc.Encounter.String() }

	switch args[0] {
	case "-reset":
		new := encounter.NewEncounter()
		*vc.Encounter = new
		return "encounter reset"

	case "add":       return e.AddMonsters(vc, args[1:])
	case "join":      return e.AddPlayer(vc)
	case "roster":    return vc.Encounter.String()

	default: return e.Help()
	}
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
