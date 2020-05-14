package verbs

import(
	"fmt"
)

// Character is stateless, in that the verb doesnb't have its own state; it simply operates
// on the character's state in the context
type Character struct{}

func (c Character)Help() string {
	return "[set field value] - fields are race,class,alignment,level,maxhp,currhp,str,..."
}

func (c Character)Process(vc VerbContext, args []string) string {
	if vc.User == "" {
		return "who are you, eh ?"
	}

	if len(args) == 0 {
		return fmt.Sprintf("%s", vc.Character)
	}
	
	switch args[0] {
	case "set":
		if len(args) != 3 { return "`set field value`, plz" }
		return vc.Character.Set(args[1], args[2])
		//c := vc.Character
		//str := c.Set(args[1], args[2])
		//vc.Character = c
		//return str

	default: return "wat?"
	}

	return ""
}
