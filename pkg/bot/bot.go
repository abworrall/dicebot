package bot

import(
	"fmt"
	"regexp"
	"strings"
	"strconv"
	"math/rand"
)

// Looks for "... dicebot <verb> <object> [<object2> ...]
func Process(user, s string) string {
	// Sanitize: remove punctuation, trim space, lowercase
	s2 := strings.ToLower(regexp.MustCompile(`([^-_a-zA-Z0-9 ])`).ReplaceAllString(s, ""))
	
	w := strings.Fields(s2)
	for len(w) > 1 {
		if w[0] != "dicebot" {
			w = w[1:]
			continue
		}
		verb,w := w[1],w[2:]

		switch verb {
		case "roll": return Roll(user, w)
		default: return "I don't know how to `"+verb+"` :("
		}
	}

	return ""
}

func Roll(user string, w []string) string {
	bits := regexp.MustCompile(`^(\d*)d(\d+)$`).FindStringSubmatch(w[0]) // 4d6, 4, 6
	if len(bits) != 3 {
		return "did not understand `"+w[0]+"` :("
	}

	n,_ := strconv.Atoi(bits[1])
	if n<=0 { n = 1 }
	if n>100 { return "way too many dice" }

	ord,_ := strconv.Atoi(bits[2])
	if ord<=0 { ord = 1 }
	if ord>100 { return "that dice is way too big" }

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
		str += fmt.Sprintf("[%s], total: %d", strings.Join(results, ","), total)
	}
	return str
}
