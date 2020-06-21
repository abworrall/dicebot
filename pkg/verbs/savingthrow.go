package verbs

import(
	"fmt"
	"strconv"

	"github.com/abworrall/dicebot/pkg/dnd5e/roll"
)

type SavingThrow struct{}

var attrModifiers = []int{
	0,-5,-4,-4,  // attr scores 0-3
	-3,-3,-2,-2,-1,-1, // attr scores 4-9
	0,0,1,1,2,2, // attr scores 10-15
	3,3,4,4,5,5, // attr scores 16-21
	6,6,7,7,8,8,9,9,10, // attr scores 22-30
}
	
func (st SavingThrow)Help() string {
	return "vs {str,int,wis,con,cha,dex,per} [DC]"
}

func (st SavingThrow)Process(vc VerbContext, args []string) string {
	if vc.Character == nil { return "who are you, again ?" }
	modifier := 0  // What we add, or remove, from the rolled die
	attrVal := 0   // The character's attribute value, the 3d6 thing
	
	if len(args) <= 1 || args[0] != "vs" { return st.Help() }
	kind,args := args[1], args[2:] // ignore "vs", shift off the attr name

	if attrVal,_ = vc.Character.Get(kind); attrVal < 0 {
		return fmt.Sprintf("you can't save against '%s'", kind)
	} else {
		modifier = attrModifiers[attrVal]
	}

	r := roll.Roll{NumDice:1, DiceSize:20, Modifier:modifier, Reason:"save vs "+kind}
	
	if len(args) == 1 {
		r.Target,_ = strconv.Atoi(args[0])
	}

	str := r.Do().String()
	vc.LogEvent(str)

	return str
}
