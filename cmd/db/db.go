package main

import(
	"bufio"
	"fmt"
	"os"

	"github.com/abworrall/dicebot/pkg/bot"
	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
	"github.com/abworrall/dicebot/pkg/verbs"
)

func main() {
	rules.Init("./data/")

	b := bot.New("dicebot", "db")
	vc := verbs.VerbContext{
		ExternalUserId: "ABCDEF123456",
		StateManager: FileStateManager{"/home/abw/db-state"},
	}
	reader := bufio.NewReader(os.Stdin)

	
	for {
		fmt.Print("> ")
		in, _ := reader.ReadString('\n')

		if out := b.ProcessLine(vc,in); out != "" {
			fmt.Printf("%s\n", out)
		}
	}
}
