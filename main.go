package main

import (
	"crypto/tls"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/thoj/go-ircevent"
	"log"
)

type ircBot struct {
	Nickname      string // The nick the bot will use
	Username      string // The IRC username the bot will use
	Password      string // Server password
	Host          string // The IRC Server hostname
	Port          int    // The IRC Server port
	UseTLS        bool   // Should we connect using TLS?
	TLSServerName string // Must supply if above is true
	Debug         bool   // Log all IRC comms to std out
}

type telegramBot struct {
	Token string
}

var ircConn *irc.Connection

func setUpIRCConnection(ircBot *ircBot) {
	ircConn = irc.IRC(ircBot.Username, ircBot.Nickname)
	ircConn.Password = ircBot.Password
	ircConn.UseTLS = ircBot.UseTLS
	ircConn.TLSConfig = &tls.Config{
		ServerName: ircBot.Host,
	}

	ircConn.VerboseCallbackHandler = ircBot.Debug

}

func main() {
	TOKEN := ""

	bot, error := tgbotapi.NewBotAPI(TOKEN)

	if error != nil {
		log.Panic(error)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	updates, error := bot.GetUpdatesChan(u)

	if error != nil {
		log.Panic(error)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// bot.Send(msg)
	}

}
