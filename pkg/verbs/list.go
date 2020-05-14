package verbs

/*

import(
	"fmt"
	"strconv"
	"strings"
	"time"
)

type List struct {
	Lists map[string][]Item
}

type Item struct {
	Time time.Time
	Description string
}

func (a Item)String() string {
	return a.Description
}
	
func (l *List)Help() string { return "[stash item] [list] [remove N] [use N]" }

func (l *List)Process(vc VerbContext, args []string) string {
	if len(args) < 1 { return l.Help() }
	
	if l.Lists == nil {
		l.Lists = map[string][]Item{}
	}

	p := LoadParty(vc)
	if vc.User == "" {
		return "no you won't, I don't know you"
	}
	user := p.UserIds[vc.User]
	if user == "" {
		return "you don't have a character; `db party claim NAME` one"
	}
	if l.Lists[user] == nil {
		l.Lists[user] = []Item{}
	}

	getN := func(args []string) (int, string) {
		if len(args) != 2 { l.Help() }
		if n,err := strconv.Atoi(args[1]); err != nil {
			return -1, fmt.Sprintf("'%s' is such nonsense", args[1])
		} else if n < 1 {
			return -1, "yes yes, very clever"
		} else if n > len(l.Lists[user]) {
			return -1, fmt.Sprintf("you don't even have %d items, let alone %d", len(l.Lists[user])+1, n)
		} else {
			return n-1, "" // -1 'cos consumers will want a slice index
		}
	}
	
	switch args[0] {
	case "-flush":
		l.Lists = nil
		
	case "stash":
		if len(args) == 1 { return "what do you want to stash, eh ?" }
		new := strings.Join(args[1:], " ")
		l.Lists[user] = append (l.Lists[user], Item{time.Now(), new})
		return "item stashed"

	case "list":
		str := ""
		for j,item := range l.Lists[user] {
			str += fmt.Sprintf("[%02d] %s\n", j+1, item)
		}
		return str

	case "use":
		if n,str := getN(args); str != "" {
			return str
		} else {
			return fmt.Sprintf("%s uses their %s\n", user, l.Lists[user][n])
		}

	case "remove":
		if n,str := getN(args); str != "" {
			return str
		} else {
			a := l.Lists[user]
			a = append(a[:n], a[n+1:]...)
			l.Lists[user] = a
		}

	default:
		return l.Help()
	}

	return ""
}

*/
