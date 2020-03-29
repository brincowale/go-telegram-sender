package telegram

import (
	"errors"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
	"net/http"
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
	if gjson.Get(body, "ok").Bool() {
		return nil
	}
	err := gjson.Get(body, "error_code").String() + ": " + gjson.Get(body, "description").String()
	return errors.New(err)
}
