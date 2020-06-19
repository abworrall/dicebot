package verbs

import(
	"fmt"
	"regexp"
	"strings"
	"strconv"
	"math/rand"
)

type Roll struct{}

func (r Roll)Help() string {
	return "4d6"
}

func (r Roll)Process(vc VerbContext, args []string) string {
	if len(args) < 1 {
		return "what do you want me to roll ?"
	}

	bits := regexp.MustCompile(`^(\d*)d(\d+)([-+]\d+)?$`).FindStringSubmatch(args[0]) // 4d6, 4, 6
	if len(bits) < 3 || len(bits) > 4 {
		return "that's not how I roll: `"+args[0]+"` is nonsense"
	}

	n,_ := strconv.Atoi(bits[1])
	if n<=0 { n = 1 }
	if n>100 { return "that's not how I roll: way too many dice" }

	ord,_ := strconv.Atoi(bits[2])
	if ord<=0 { ord = 1 }
	if ord>100 { return "that's not how I roll: that dice is way too big" }

	modifier := 0
	if len(bits) == 4 {
		modifier,_ = strconv.Atoi(bits[3])
	}

	results := []string{}
	total := 0
	for i:=0; i<n; i++ {
		r := rand.Intn(ord) + 1 // returns a value in the range [0,ord)
		results = append(results, fmt.Sprintf("%d", r))
		total += r
	}

	modstr := ""
	if modifier != 0 {
		total += modifier
		modstr = fmt.Sprintf("%+d", modifier)
	}

	str := fmt.Sprintf("%dd%d%s: you got ", n, ord, modstr)
	str += fmt.Sprintf("[%s] %s total:%d", strings.Join(results, ","), modstr, total)

	if len(args) > 1 {
		str += " (for " + strings.Join(args[1:], " ") + ")"
	}

	vc.LogEvent("rolled " + str)

	return str
}
