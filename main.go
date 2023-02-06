package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	dd = "Bot: 我是哈璐，是一隻狗狗，英文名字叫Haru，我的主人是我自己"
	br = "\n"

	// maximumContextLength int64 = 4097
	maxTokens      int64 = 900
	tokenThreshold int64 = 3000
)

var (
	//ps    string = dd
	lock = new(sync.RWMutex)
	ps   = []string{dd}
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
					case message.Text == "/reset":
						reset()
					case strings.HasPrefix(message.Text, "/t "):
						msg := strings.Replace(message.Text, "/t ", "", 1)
						if msg == "" {
							break
						}
						push(youSay(msg))
						var p string
						var totalTokens int64
						p = prompt() + botSay("")

						response, totalTokens, err = generateChatResponse(
							Chat{
								Model:             "text-davinci-003",
								Prompt:            p,
								N:                 1,
								MaxTokens:         maxTokens,
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
						if strings.HasPrefix(reply, br) {
							reply = strings.Replace(reply, br, "", 1)
						}

						push(botSay(reply))

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
							log.Print(err)
						}

						if totalTokens > tokenThreshold {
							cut()
							log.Print("do cut, total_tokens:", totalTokens)
						}
					case strings.HasPrefix(message.Text, "/x "):
						msg := strings.Replace(message.Text, "/x ", "", 1)
						push(botSay(msg))
					case strings.HasPrefix(message.Text, "/e "):
						msg := strings.Replace(message.Text, "/e ", "", 1)
						push(msg)
					case strings.HasPrefix(message.Text, "/ti "):
						msg := strings.Replace(message.Text, "/ti ", "", 1)
						if msg == "" {
							break
						}
						push(youSay(msg))
						var p string
						var totalTokens int64
						p = promptLast() + botSay("")

						response, totalTokens, err = generateChatResponse(
							Chat{
								Model:             "text-davinci-003",
								Prompt:            p,
								N:                 1,
								MaxTokens:         maxTokens,
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
						if strings.HasPrefix(reply, br) {
							reply = strings.Replace(reply, br, "", 1)
						}

						push(botSay(reply))

						if totalTokens > tokenThreshold {
							cut()
							log.Print("do cut, total_tokens:", totalTokens)
						}

						response, err = generateImageResponse(
							Image{
								Prompt: prompt(),
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

func say(who, m string) (msg string) {
	return who + ": " + m
}

func botSay(m string) (msg string) {
	return say("Bot", m)
}

func youSay(m string) (msg string) {
	return say("You", m)
}

func promptLast() (prompt string) {
	lock.RLock()
	g := 3
	if len(ps) < g {
		g = len(ps)
	}
	prompt = strings.Join(ps[len(ps)-g:], br) + br
	lock.RUnlock()
	return
}

func prompt() (prompt string) {
	lock.RLock()
	prompt = strings.Join(ps, br) + br
	lock.RUnlock()
	return
}

func push(msg string) {
	lock.Lock()
	ps = append(ps, msg)
	lock.Unlock()
}

func reset() {
	lock.Lock()
	ps = []string{dd}
	lock.Unlock()
}

func cut() {
	lock.Lock()
	ps = append(ps[0:1], ps[5:]...)
	lock.Unlock()
}
