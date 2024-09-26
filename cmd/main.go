package main

import (
	"log"

	bot "github.com/YotoHana/tgbot-golang/internal/app"
)

func main()  {
	err := bot.Run()
	if err != nil {
		log.Fatalf("Error with bot: %v", err)
	}
}