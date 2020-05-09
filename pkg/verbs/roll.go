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

func (r Roll)Process(c VerbContext, args []string) string {
	if len(args) < 1 {
		return "what do you want me to roll ?"
	}

	bits := regexp.MustCompile(`^(\d*)d(\d+)$`).FindStringSubmatch(args[0]) // 4d6, 4, 6
	if len(bits) != 3 {
		return "that's not how I roll: `"+args[0]+"` is nonsense"
	}

	n,_ := strconv.Atoi(bits[1])
	if n<=0 { n = 1 }
	if n>100 { return "that's not how I roll: way too many dice" }

	ord,_ := strconv.Atoi(bits[2])
	if ord<=0 { ord = 1 }
	if ord>100 { return "that's not how I roll: that dice is way too big" }

	results := []string{}
	total := 0
	for i:=0; i<n; i++ {
		r := rand.Intn(ord) + 1 // returns a value in the range [0,ord)
		results = append(results, fmt.Sprintf("%d", r))
		total += r
	}

	str := fmt.Sprintf("%dd%d: you got ", n, ord)

	if n == 1 {
		str += results[0]
	} else {
		str += fmt.Sprintf("[%s]  total:%d", strings.Join(results, ","), total)
	}
	return str
}
