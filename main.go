package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {

				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					// Use OpenAI API to generate response
					var response string
					switch {
					case strings.HasPrefix(message.Text, "/t "):
						response, err = generateChatResponse(
							Chat{
								Model:     "text-davinci-003",
								Prompt:    strings.Replace(message.Text, "/t ", "", 1),
								N:         1,
								MaxTokens: 2048,
							})
						if err != nil {
							log.Print(err)
							break
						}
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
							log.Print(err)
						}
					case strings.HasPrefix(message.Text, "/i "):
						response, err = generateImageResponse(
							Image{
								Prompt: strings.Replace(message.Text, "/i ", "", 1),
								N:      1,
								Size:   "256x256",
							})
						if err != nil {
							log.Print(err)
							break
						}
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(response, response)).Do(); err != nil {
							log.Print(err)
						}
					default:
						return
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy, or something else.
	if err := http.ListenAndServe(":1337", nil); err != nil {
		log.Fatal(err)
	}
}
