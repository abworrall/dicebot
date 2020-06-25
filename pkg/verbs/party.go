package verbs

import(
	"sort"

	"github.com/abworrall/dicebot/pkg/character"
)

// Party is stateful - it records a simple list of character names.
// Various subverbs iterate over all party members.
type Party struct{
	Names map[string]bool
}

func (p *Party)Help() string { return "[add NAME], [remove NAME], [rest]" }

func (p *Party)Process(vc VerbContext, args []string) string {
	if p.Names == nil {
		p.Names = map[string]bool{}
	}

	// Handle changes to the party list and bail
	if len(args) == 2 {
		switch args[0] {
		case "add":
			p.Names[args[1]] = true
			return "added"
		case "delete":
			delete(p.Names, args[1])
			return "deleted"
		}
	}

	// Decide which action to take
	f := CharOneliner
	if len(args) > 0 {
		switch args[0] {
		case "rest": f = CharLongRest
		}
	}

	chars := make([]*character.Character, len(p.Names))
	for i,name := range p.ListNames() {
		chars[i] = vc.loadCharacter(name)
	}

	str := ""
	for _,c := range chars {
		str += f(c)
	}

	for _,c := range chars {
		vc.maybeSaveCharacter(c)
	}

	return str
}

type CharFunc func(*character.Character) string

func CharOneliner(c *character.Character) string {
	return c.Summary() + "\n"
}

func CharLongRest(c *character.Character) string {
	// Restore the character as per a long rest
	c.CurrHitpoints = c.MaxHitpoints
	c.Slots.Reset()
	return CharOneliner(c)
}

func (p *Party)ListNames() []string {
	names := []string{}
	for k,_ := range p.Names {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
