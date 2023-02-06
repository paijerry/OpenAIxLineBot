package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/paijerry/ezapi"
)

type Chat struct {
	Model             string  `json:"model"`
	Prompt            string  `json:"prompt"`
	N                 int64   `json:"n"`
	Temperature       float64 `json:"temperature"`
	TopP              int64   `json:"top_p"`
	MaxTokens         int64   `json:"max_tokens"`
	Frequency_Penalty float64 `json:"frequency_penalty"`
	Presence_Penalty  float64 `json:"presence_penalty"`
}

type RspnChat struct {
	Createed int64          `json:"created"`
	Choices  []RspnChatData `json:"choices"`
	Usage    RspnUsage      `json:"usage"`
}

type RspnChatData struct {
	Text string `json:"text"`
}

type RspnUsage struct {
	TotalTokens int64 `json:"total_tokens"`
}

func generateChatResponse(req Chat) (o string, totalTokens int64, err error) {

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", "Bearer "+os.Getenv("OpenAiToken"))
	url := "https://api.openai.com/v1/completions"
	reqByte, err := json.Marshal(req)
	if err != nil {
		return
	}

	rspn, err := ezapi.New().URL(url).Header(header).JSON(reqByte).TimeOut(60).Do("POST")
	if err != nil {
		return
	}
	if rspn.StatusCode != 200 {
		err = errors.New("HTTPStatus:" + string(rspn.StatusCode) + " => " + string(rspn.Body))
		return
	}
	fmt.Println(string(rspn.Body))
	result := RspnChat{}
	err = json.Unmarshal(rspn.Body, &result)
	if err != nil {
		err = errors.New("Error:" + err.Error() + " => " + string(rspn.Body))
		return
	}

	if len(result.Choices) == 0 {
		err = errors.New("Error: No data => " + string(rspn.Body))
		return
	}
	//fmt.Println(result.Choices[0].Text)

	// if strings.HasPrefix(result.Choices[0].Text, "\n\n") {
	// 	result.Choices[0].Text = strings.Replace(result.Choices[0].Text, "\n\n", "", 1)
	// }

	return result.Choices[0].Text, result.Usage.TotalTokens, nil
}
