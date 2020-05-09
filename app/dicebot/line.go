package main

import(
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	linebothttphandler "github.com/line/line-bot-sdk-go/linebot/httphandler"
	
	mybot "github.com/abworrall/dicebot/pkg/bot"
	"github.com/abworrall/dicebot/pkg/config"
)

// https://github.com/line/line-bot-sdk-go/blob/master/examples/echo_bot_handler/server.go
func registerLineHandlerFor(url string) {
	handler, err := linebothttphandler.New(
		config.Get("line.channelsecret"),
		config.Get("line.channeltoken"),
	)
	if err != nil {
		log.Fatal(err)
	}

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

			user := ""
			if ev.Source.Type == linebot.EventSourceTypeUser {
				user = ev.Source.UserID
			}

			switch msg := ev.Message.(type) {
			case *linebot.TextMessage:
				if respText := mybot.Process(user, msg.Text); respText != "" {
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
