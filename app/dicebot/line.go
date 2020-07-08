package main

import(
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	linebothttphandler "github.com/line/line-bot-sdk-go/linebot/httphandler"

	"github.com/abworrall/dicebot/pkg/config"
	mybot "github.com/abworrall/dicebot/pkg/bot"
	"github.com/abworrall/dicebot/pkg/state"
	"github.com/abworrall/dicebot/pkg/verbs"
)

// https://github.com/line/line-bot-sdk-go/blob/master/examples/echo_bot_handler/server.go
func registerLineHandlerFor(url string, gcpProjectId string) {
	handler, err := linebothttphandler.New(
		config.Get("line.channelsecret"),
		config.Get("line.channeltoken"),
	)
	if err != nil {
		log.Fatal(err)
	}

	db := mybot.New("dicebot", "db")

	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			log.Printf("handler.NewClient: %v", err)
			return
		}

		for _,ev := range events {
			// https://developers.line.biz/en/docs/messaging-api/receiving-messages/#webhook-event-types
			// https://github.com/line/line-bot-sdk-go/blob/master/linebot/event.go
			if ev.Type != linebot.EventTypeMessage { continue }

			ctx := r.Context()
			vc := verbs.VerbContext{
				Ctx: ctx,
				StateManager: state.NewGcpStateManager(ctx, gcpProjectId),
			}
			switch ev.Source.Type {
			case linebot.EventSourceTypeUser:  vc.ExternalUserId = ev.Source.UserID
			case linebot.EventSourceTypeGroup: vc.ExternalUserId = ev.Source.UserID
			}
			vc.Debug = fmt.Sprintf("ev<%T>: %#v\n\nsrc<%T>: %#v\n\n", ev, ev, *ev.Source, *ev.Source)

			switch msg := ev.Message.(type) {
			case *linebot.TextMessage:
				vc.Debug += fmt.Sprintf("msg<%T>: %#v\n\n", msg, msg)

				if respText := db.ProcessLine(vc, msg.Text); respText != "" {
					if _,err := bot.ReplyMessage(ev.ReplyToken, linebot.NewTextMessage(respText)).Do(); err != nil {
						log.Printf("bot.ReplyMessage: %v", err)
					}
				}
			default:
				log.Printf("event was confusing: %s", ev)
			}
		}
	})

	http.Handle(url, handler)
}
