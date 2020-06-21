package verbs

import(
	"strings"

	"github.com/abworrall/dicebot/pkg/dnd5e/roll"
)

type Roll struct{}

func (r Roll)Help() string {
	return "4d6+3 >=8 withadvantage"
}

func (r Roll)Process(vc VerbContext, args []string) string {
	if len(args) < 1 {
		return "what do you want me to roll ?"
	}
	
	str := roll.New(strings.Join(args, " ")).String()
	vc.LogEvent("rolled " + str)
	return str
}
