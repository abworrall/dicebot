package verbs

import(
	"fmt"
	"sort"
	"strconv"

	"github.com/abworrall/dicebot/pkg/character"
)

type Party struct{
	Name string // party name
	Characters map[string]character.Character // key is name

	UserIds map[string]string // Map some kind of user identifier to a character name
}

func (p *Party)Help() string {
	return "[claim NAME] [delete NAME] [add NAME stats...]"
}

func (p *Party)String() string {
	keys := []string{}
	for k,_ := range p.Characters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	str := ""
	for _,k := range keys {
		str += fmt.Sprintf("%s\n", p.Characters[k])
	}
	return str
}

// party                                          - print out all the characters
// party claim name                               - associates current user with the named char

// party add    name str int wis con cha dex per  - adds/updates a character
// party delete name                              - removes a character

func (p *Party)Process(vc VerbContext, args []string) string {
	if len(args) == 0 {
		return fmt.Sprintf("%s", p)
	}

	switch args[0] {
	case "add": return p.Add(args[1:])
	case "delete": return p.Delete(args[1])
	case "claim": return p.Claim(vc, args[1])
	default: return "party's over"
	}

	return ""
}

func (p *Party)Claim(vc VerbContext, name string) string {
	if p.UserIds == nil {
		p.UserIds = map[string]string{}
	}

	p.UserIds[vc.User] = name
	return fmt.Sprintf("%s has been claimed by %s", name, vc.User)
}

func (p *Party)Delete(name string) string {
	delete(p.Characters, name)
	return fmt.Sprintf("%s has left the party", name)
}

func (p *Party)Add(args []string) string {
	if len(args) != 8 {
		return "can't add that, need [name str int wis con cha dex per]"
	}

	myAtoi := func(s string) int {
		if v,err := strconv.Atoi(s); err != nil {
			return 0
		} else {
			return v
		}
	}

	c := character.Character{
		Name: args[0],
		Str: myAtoi(args[1]),
		Int: myAtoi(args[2]),
		Wis: myAtoi(args[3]),
		Con: myAtoi(args[4]),
		Cha: myAtoi(args[5]),
		Dex: myAtoi(args[6]),
		Per: myAtoi(args[7]),
	}

	if p.Characters == nil {
		p.Characters = map[string]character.Character{}
	}
	p.Characters[c.Name] = c
	
	return fmt.Sprintf("%s has joined the party", c.Name)
}
