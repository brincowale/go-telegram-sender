package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-retryablehttp"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	Request *http.Client
	Token   string
}

type Message struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func New(token string) Client {
	httpRetryClient := retryablehttp.NewClient()
	httpRetryClient.RetryMax = 3
	httpRetryClient.HTTPClient.Timeout = time.Second * 60
	httpClient := httpRetryClient.StandardClient()
	return Client{
		Request: httpClient,
		Token:   token,
	}
}

func (c Client) SendMessage(chatId string, message string) error {
	data, err := json.Marshal(
		Message{
			ChatId: chatId,
			Text:   message,
		},
	)
	if err != nil {
		return err
	}
	endpoint := "https://api.telegram.org/bot" + c.Token + "/sendMessage?parse_mode=HTML"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Request.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var response TelegramResponse
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return err
	}
	if response.Ok {
		return nil
	}
	strError := strconv.Itoa(response.ErrorCode) + ": " + response.Description
	return errors.New(strError)
}
