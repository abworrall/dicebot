package verbs

import "fmt"

// Character is stateless, in that the verb doesn't have its own state; it simply operates
// on the character's state in the context
type Character struct{}

func (c Character)Help() string {
	return "[set FIELD VALUE], [list]"
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

	case "list":
		return "[Useful fields: weapon, armor, shield]\n"+
			"[Less useful fields: race, class, alignment, level, maxhp, buff, hp, str, int, wis, con, cha, dex per]"

	case "setstats":
		if len(args) != 8 { return "`setstats 1 2 3 4 5 6 7`, plz" }
		vc.Character.Set("str", args[1])
		vc.Character.Set("int", args[2])
		vc.Character.Set("wis", args[3])
		vc.Character.Set("con", args[4])
		vc.Character.Set("cha", args[5])
		vc.Character.Set("dex", args[6])
		vc.Character.Set("per", args[7])
		
	default: return "wat?"
	}

	return ""
}
