package telegram

import (
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	Request *gorequest.SuperAgent
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
	Result      struct {
		MessageID int `json:"message_id"`
		Chat      struct {
			ID       int64  `json:"id"`
			Title    string `json:"title"`
			Username string `json:"username"`
			Type     string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}

func New(token string) Client {
	r := gorequest.New().Set("Content-Type", "application/json").
		Timeout(30*time.Second).Retry(3, 5*time.Second, http.StatusInternalServerError)
	return Client{
		Request: r,
		Token:   token,
	}
}

func SendMessage(c Client, m Message) error {
	var URL = "https://api.telegram.org/bot" + c.Token + "/sendMessage"
	data := Message{
		ChatId: m.ChatId,
		Text:   m.Text,
	}
	_, body, errs := c.Request.Post(URL).Send(data).End()
	if errs != nil {
		return errs[0]
	}
	var response TelegramResponse
	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		return err
	}
	if response.Ok {
		return nil
	}
	strError := strconv.Itoa(response.ErrorCode) + ": " + response.Description
	return errors.New(strError)
}
