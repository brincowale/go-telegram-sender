# go-telegram-sender
A simple project to send messages to Telegram from Go projects

Use the `go-retryablehttp` to retry automatically when Telegram returns a 429 HTTP error

# Usage

```go
package main

import (
	"fmt"
	telegram "github.com/brincowale/go-telegram-sender"
)

func main() {
	t := telegram.New("BOT-API-KEY")
	err := t.SendMessage("@channelID", "Hello World!")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK")
	}
}
```
