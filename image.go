package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/paijerry/ezapi"
)

type Image struct {
	Prompt string `json:"prompt"`
	N      int64  `json:"n"`
	Size   string `json:"size"`
}

type RspnImage struct {
	Createed int64           `json:"created"`
	Data     []RspnImageData `json:"data"`
}

type RspnImageData struct {
	Url string `json:"url"`
}

func generateImageResponse(req Image) (o string, err error) {

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", "Bearer "+os.Getenv("OpenAiToken"))
	url := "https://api.openai.com/v1/images/generations"
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

	result := RspnImage{}
	err = json.Unmarshal(rspn.Body, &result)
	if err != nil {
		err = errors.New("Error:" + err.Error() + " => " + string(rspn.Body))
		return
	}

	if len(result.Data) == 0 {
		err = errors.New("Error: No data => " + string(rspn.Body))
		return
	}

	return result.Data[0].Url, nil
}
