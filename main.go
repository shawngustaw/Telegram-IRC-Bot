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

func getToken() string {
	token := os.Getenv("TELEGRAM_TOKEN")

	if token == "" {
		fmt.Errorf("Telegram Token not set. Exiting.")
		os.Exit(1)
	}

	return token
}

func main() {
	// Telegram reads, IRC writes
	ping := make(chan string)
	// IRC reads, Telegram writes
	pong := make(chan string)

	TOKEN := getToken()

	Telegram := newTelegramBot(TOKEN, true, ping, pong)
	Telegram.initConnection()
	Telegram.beginLoop()

	IRC := newIRCBot(
		"Abot", "Abot", "irc.oftc.net", 6667, "#tag",
		true, pong, ping)
	IRC.initConnection()
	IRC.beginLoop()

	for {
		// keep main program running
	}

	return

}
