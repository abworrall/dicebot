package bot

import(
	"regexp"
	"strings"

	"github.com/abworrall/dicebot/pkg/verbs"
)

type Bot struct{
	Name string // How the bot will refer to itself
	nicknames map[string]bool // How the bot knows it is being invoked
}

func New(name string, nicknames ...string) Bot {
	b := Bot{
		Name: name,
		nicknames: map[string]bool{name:true},
	}

	for _,nick := range nicknames {
		b.nicknames[strings.ToLower(nick)] = true
	}

	return b
}

func (b Bot)ProcessLine(vc verbs.VerbContext, in string) string {
	w := sanitizeLine(in)

	if _,exists := b.nicknames[w[0]]; !exists { return "" }
	if len(w) < 2 { return "Yo" }

	return verbs.Act(vc, w[1], w[2:])
}


func sanitizeLine(in string) ([]string) {
	// Sanitize: remove punctuation, trim space, lowercase. Leave in colons for spell indices
	sanitized := strings.ToLower(regexp.MustCompile(`([^-+:_a-zA-Z0-9 ])`).ReplaceAllString(in, ""))

	return strings.Fields(sanitized) // Also trims space
}
