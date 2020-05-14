package verbs

import(
	"math/rand"
	"regexp"
	"strings"
)

type Insult struct{
	Insults []string
}

func (i *Insult)Help() string { return "[user] [learn user is a limpet]" }

func (i *Insult)Process(vc VerbContext, args []string) string {
	if len(args) == 0 { return i.Help() }

	if i.Insults == nil {
		i.Insults = []string{}
	}
	
	switch args[0] {
	case "-flush":
		i.Insults = []string{}
		
	case "learn":
		new := strings.Join(args[1:], " ")
		i.Insults = append(i.Insults, new)
		return "heh heh"
		
	default:
		str := i.Insults[rand.Intn(len(i.Insults))]
		target := strings.Join(args, " ")
		str = regexp.MustCompile(`\buser\b`).ReplaceAllString(str, target)
		return str
	}

	return ""
}
