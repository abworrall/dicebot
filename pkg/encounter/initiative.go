package encounter

import(
	"fmt"
	"sort"
	"strings"
)

type Initiative struct {
	Scores map[int][]string
}

func NewInitiative() Initiative {
	return Initiative{
		Scores: map[int][]string{},
	}
}

func (init *Initiative)String() string {
	vals := []int{}
	for val,_ := range init.Scores {
		vals = append(vals, val)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(vals)))

	strs := []string{}
	for _,val := range vals {
		names := init.Scores[val]
		strs = append(strs, fmt.Sprintf("%d:%s", val, strings.Join(names, ",")))
	}
	
	str := fmt.Sprintf("{Init: %s}", strings.Join(strs, "; "))

	return str
}

func (init *Initiative)Set(name string, val int) {
	if init.Scores[val] == nil {
		init.Scores[val] = []string{}
	}
	init.Scores[val] = append(init.Scores[val], name)
}

func (init *Initiative)Get(s string) int{
	for val, names := range init.Scores {
		for _,name := range names {
			if s == name {
				return val
			}
		}
	}
	return -1
}
