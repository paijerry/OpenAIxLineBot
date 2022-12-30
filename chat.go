package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/paijerry/ezapi"
)

type Chat struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	N         int64  `json:"n"`
	MaxTokens int64  `json:"max_tokens"`
}

type RspnChat struct {
	Createed int64          `json:"created"`
	Choices  []RspnChatData `json:"choices"`
}

type RspnChatData struct {
	Text string `json:"text"`
}

func generateChatResponse(req Chat) (o string, err error) {

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", "Bearer "+os.Getenv("OpenAiToken"))
	url := "https://api.openai.com/v1/completions"
	reqByte, err := json.Marshal(req)
	if err != nil {
		return
	}

	rspn, err := ezapi.New().URL(url).Header(header).JSON(reqByte).Do("POST")
	if err != nil {
		return
	}
	if rspn.StatusCode != 200 {
		err = errors.New("HTTPStatus:" + string(rspn.StatusCode) + " => " + string(rspn.Body))
		return
	}

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
	fmt.Println(result.Choices[0].Text)
	if strings.HasPrefix(result.Choices[0].Text, "\n\n") {
		result.Choices[0].Text = strings.Replace(result.Choices[0].Text, "\n\n", "", 1)
	}

	return result.Choices[0].Text, nil
}
