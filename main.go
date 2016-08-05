package main

import (
	//"crypto/tls"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/thoj/go-ircevent"
	"time"
)

type message struct {
	Sender string
	Text   string
	Sent   time.Time
}

var ircConn *irc.Connection
var BotAPI *tgbotapi.BotAPI

func main() {
	// Telegram reads, IRC writes
	ping := make(chan string)
	// IRC reads, Telegram writes
	pong := make(chan string)

	setUpTelegramConnection(ping, pong)
	IRC := newIRCBot(
		"Abot", "Abot", "irc.oftc.net", 6667, "#tag",
		true, pong, ping)
	IRC.initConnection()
	IRC.beginLoop()

	for {
		// keep main program running
	}

	return // BEGIN TELEGRAM CODE BELOW

}
