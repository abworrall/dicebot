package main

import(
	"bufio"
	"fmt"
	"os"

	"github.com/abworrall/dicebot/pkg/bot"
	"github.com/abworrall/dicebot/pkg/dnd5e/rules"
	"github.com/abworrall/dicebot/pkg/verbs"
	"github.com/abworrall/dicebot/pkg/state"
)

func main() {
	rules.Init("./data/")

	b := bot.New("dicebot", "db")
	vc := verbs.VerbContext{
		ExternalUserId: "ABCDEF123456",
		StateManager: state.FileStateManager{"/home/abw/db-state"},
	}

/* Sadly, this all dies :(

panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x8fb180]

goroutine 1 [running]:
go.opencensus.io/trace.FromContext(...)
	/home/abw/go/pkg/mod/go.opencensus.io@v0.22.4/trace/trace.go:113
go.opencensus.io/trace.StartSpan(0x0, 0x0, 0xb1f1ed, 0x21, 0x0, 0x0, 0x0, 0x106c458, 0x203000, 0x203000)
	/home/abw/go/pkg/mod/go.opencensus.io@v0.22.4/trace/trace.go:172 +0x70
	
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" && os.Getenv("GOOGLE_CLOUD_PROJECT") != "" {
		gcpProjectId :=	os.Getenv("GOOGLE_CLOUD_PROJECT")
		ctx := trace.NewContext(context.Background(), nil)
		//ctx := context.Background()
		vc.StateManager = state.NewGcpStateManager(ctx, gcpProjectId)
		fmt.Printf("(connecting to GCP !)\n")
	}

*/

	reader := bufio.NewReader(os.Stdin)

	
	for {
		fmt.Print("> ")
		in, _ := reader.ReadString('\n')

		if out := b.ProcessLine(vc,in); out != "" {
			fmt.Printf("%s\n", out)
		}
	}
}
