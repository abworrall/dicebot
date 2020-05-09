package main

import(
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/abworrall/dicebot/pkg/bot"
	"github.com/abworrall/dicebot/pkg/verbs"
)

func init() {
	http.HandleFunc("/debug", debugHandler)
	registerLineHandlerFor("/line", os.Getenv("GOOGLE_CLOUD_PROJECT"))
}


func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[dicebot] listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func req2ctx(r *http.Request) context.Context {
	ctx,_ := context.WithTimeout(r.Context(), 55 * time.Second)
	return ctx
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	vc := verbs.VerbContext{User: "DEADBEEF"}
	b := bot.New("dicebot", "db")
	str := b.ProcessLine(vc, r.FormValue("q"))

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("OK!\ndicebot debug handler\nin  [%s]\nout [%s]\n", r.FormValue("q"), str)))
}
