package verbs

import(
	"fmt"
	"math/rand"
)

type SavingThrow struct{}

var modifiers = map[string]int{
	"veryeasy": -8,
	"easy": -4,
	"hard": 4,
	"veryhard": 8,
}

func (st SavingThrow)Help() string {
	return "vs {str,int,wis,con,cha,dex} [{veryeasy, easy, hard, veryhard}]"
}

func (st SavingThrow)Process(vc VerbContext, args []string) string {
	if vc.Character == nil { return "who are you, again ?" }
	modifier := 0 // What we add, or remove, from the rolled die
	val := 0      // The character's attribute value, the 3d6 thing
	
	if len(args) == 0 || args[0] != "vs" { return st.Help() }
	args = args[1:] // shift

	if len(args) == 0 { return st.Help() }
	kind,args := args[0], args[1:]

	if len(args) == 1 {
		if m,exists := modifiers[args[0]]; !exists {
			return fmt.Sprintf("I don't have a '%s' roll modifier ", args[0])
		} else {
			modifier = m
		}
	}

	if val,_ = vc.Character.Get(kind); val < 0 {
		return fmt.Sprintf("you can't save against '%s'", kind)
	}

	x := rand.Intn(20) + 1 // d20
	madeSave := (x + modifier) <= val
	outcome := "SAVE!"
	if !madeSave { outcome = "failed :(" }
	
	modStr := ""
	if modifier != 0 {
		modStr = fmt.Sprintf(" with modifier %d", modifier)
	}

	vc.LogEvent(fmt.Sprintf("saved vs %s[%2d]%s: got %2d, %s", kind, val, modStr, x, outcome))

	str := fmt.Sprintf("%s, save vs %s[%2d]%s: you rolled %2d, %s", vc.User, kind, val, modStr, x, outcome)

	return str
}
