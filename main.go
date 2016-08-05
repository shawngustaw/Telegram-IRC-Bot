package main

import (
	//"crypto/tls"
	"fmt"
	"os"
	"time"
)

type message struct {
	Sender string
	Text   string
	Sent   time.Time
}

func getEnvVariable(name string) string {
	token := os.Getenv(name)

	if token == "" {
		fmt.Printf("%s Token not set. Exiting.", name)
		os.Exit(1)
	}

	return token

}

func getTelegramToken() string {
	return getEnvVariable("TELEGRAM_TOKEN")
}

func getLCBToken() string {
	return getEnvVariable("LCB_TOKEN")
}

func main() {
	// Telegram reads, IRC writes
	ping := make(chan string)
	// IRC reads, Telegram writes
	pong := make(chan string)

	TelegramToken := getTelegramToken()

	Telegram := newTelegramBot(
		TelegramToken,
		true,
		ping,
		pong)
	Telegram.initConnection()
	Telegram.beginLoop()

	LCBToken := getLCBToken()
	DevelopmentRoom := "5047c30359e957b86a000001"

	LetsChat := newLetsChatBot(
		LCBToken,
		"cairo.sdelements.com",
		DevelopmentRoom,
		pong,
		ping)

	LetsChat.beginLoop()

	/*
		IRC := newIRCBot(
			"Abot",
			"Abot",
			"irc.oftc.net",
			 6667,
			 "#tag",
			true,
			 pong,
			  ping)
		IRC.initConnection()
		IRC.beginLoop()
	*/

	for {
		// keep main program running
	}

	return

}
