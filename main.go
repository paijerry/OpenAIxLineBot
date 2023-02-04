package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	dd = "Bot: 我是哈璐，是一隻狗狗，英文名字叫Haru，我的主人是我自己\n"
)

var ps string = dd

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
					case message.Text == "/reset":
						ps = dd
					case strings.HasPrefix(message.Text, "/t "):
						msg := strings.Replace(message.Text, "/t ", "", 1)
						var p string
						if ps == "" {
							p = msg + "\n"
						} else {
							p = ps + "You: " + msg + "\nBot:"
						}

						response, err = generateChatResponse(
							Chat{
								Model:             "text-davinci-003",
								Prompt:            p,
								N:                 1,
								MaxTokens:         2048,
								Temperature:       0.9,
								TopP:              1,
								Frequency_Penalty: 0.5,
								Presence_Penalty:  0,
							})
						if err != nil {
							log.Print(err)
							break
						}

						reply := response
						if strings.HasPrefix(reply, "\n") {
							reply = strings.Replace(reply, "\n", "", 1)
						}

						if ps == "" {
							ps = "You: " + p + "Bot:" + response + "\n"
						} else {
							ps = p + response + "\n"
						}

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
