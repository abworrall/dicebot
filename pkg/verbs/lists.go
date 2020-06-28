package verbs

import(
	"fmt"
	"sort"
	"strconv"
	"strings"
)


// Lists is a stateful verb, holding some lists for the whole group.
type Lists struct {
	Lists map[string]List
}

func (l *Lists)Help() string { return "NAME [add [NN] THING], [remove [NN] THING]" }

// list
// list quests
// list quests add kill the thing
// list quests remove kill the thing
// list cash add 15 gp
// list cash remove 4 gp

func (l *Lists)Process(vc VerbContext, args []string) string {
	if l.Lists == nil {
		l.Lists = map[string]List{}
	}

	if len(args) == 0 {
		return l.String()
	} else if len(args) == 1 {
		if _,exists := l.Lists[args[0]]; exists {
			return l.Lists[args[0]].String()
		} else {
			return fmt.Sprintf("You don't have a list of '%s'", args[0])
		}
	}

	if _,exists := l.Lists[args[0]]; !exists {
		l.Lists[args[0]] = NewList()
	}

	list, action, args := args[0], args[1], args[2:]
	n := 1

	// If there is a quantity, shift it off
	if len(args) > 1 {
		if val,err := strconv.Atoi(args[0]); err == nil {
			n = val
			args = args[1:]
		}
	}

	if action == "remove" {
		n *= -1
	} else if action != "add" {
		return l.Help()
	}

	l.Lists[list].Update(strings.Join(args, " "), n)

	return l.Lists[list].String()
}


func (l *Lists)String() string {
	keys := []string{}
	for k,_ := range l.Lists {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	str := ""
	for _,k := range keys {
		str += fmt.Sprintf("--{ %s }--\n%s\n", k, l.Lists[k])
	}
	return str
}


type List map[string]int // The value is the count of the item

func NewList() List {
	return map[string]int{}
}

func (l List)String() string {
	keys := []string{}
	for k,_ := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	str := ""
	for _,k := range keys {
		str += k
		if l[k] <= 0 {
			str += " (none left!)"
		} else if l[k] > 1 {
			str += fmt.Sprintf(" (%d)", l[k])
		}
		str += "\n"
	}
	return str
}
	
func (l List)Add(i string) {
	l.Update(i, 1)
}

func (l List)Update(i string, n int) {
	l[i] += n

	if n == -1 && l[i] == 0 {
		l.Remove(i)
	}
}

func (l List)Remove(i string) {
	delete (l, i)
}
