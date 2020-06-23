package verbs

import(
	"math/rand"
	"regexp"
	"strings"
)

type Insult struct{
	Insults []string
	New bool
}

func (i *Insult)Help() string { return "[USER], [learn USER is a limpet]" }

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
		i.New = true
		return "heh heh"
		
	default:
		if len(i.Insults) == 0 { return "teach me" }
		str := i.Insults[rand.Intn(len(i.Insults))]
		if i.New {
			// If we've just learned a new one, use it stratight away :)
			str = i.Insults[len(i.Insults)-1]
			i.New = false
		}
		target := strings.Join(args, " ")
		str = regexp.MustCompile(`\buser\b`).ReplaceAllString(str, target)
		return str
	}

	return ""
}
