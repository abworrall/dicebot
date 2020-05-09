package verbs

import(
	"fmt"
	"math/rand"
	"strconv"
)

type SavingThrow struct{}

func (st SavingThrow)Help() string {
	return "vs {str,int,wis,con,cha,dex} [as NAME] [modifier, added to roll; -ve means easier]"
}

func (st SavingThrow)Process(vc VerbContext, args []string) string {
	p := LoadParty(vc)
	user := ""
	modifier := 0 // What we add, or remove, from the rolled die
	val := 0      // The character's attribute value, the 3d6 thing
	
	if len(args) == 0 || args[0] != "vs" { return st.Help() }
	args = args[1:] // shift

	if len(args) == 0 { return st.Help() }
	kind,args := args[0], args[1:]
	
	if len(args) > 1 && (args[0] == "for" || args[0] == "as") {
		user = args[1]
		args = args[2:]
	}

	if len(args) == 1 {
		modifier,_ = strconv.Atoi(args[0])
	}

	if user == "" {
		if n,exists := p.UserIds[vc.User]; !exists {
			return "You haven't claimed a character in the party"
		} else {
			user = n
		}
	}
	if c,exists := p.Characters[user]; ! exists {
		return fmt.Sprintf("'%s' hasn't joined the party", user)
	} else {
		switch kind {
		case "str": val = c.Str
		case "int": val = c.Int
		case "wis": val = c.Wis
		case "con": val = c.Con
		case "cha": val = c.Cha
		case "dex": val = c.Dex
		case "per": val = c.Per
		default: return fmt.Sprintf("you can't save against '%s'", kind)
		}
	}

	x := rand.Intn(20) + 1 // d20
	madeSave := (x + modifier) <= val
	outcome := "SAVE!"
	if !madeSave { outcome = "failed :(" }
	
	modStr := ""
	if modifier != 0 {
		modStr = fmt.Sprintf(" with modifier %d", modifier)
	}

	str := fmt.Sprintf("%s saves vs %s[%2d]%s: you rolled %2d, %s", user, kind, val, modStr, x, outcome)
	return str
}
