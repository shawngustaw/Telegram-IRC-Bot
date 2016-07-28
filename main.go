package main

import (
	//"crypto/tls"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/thoj/go-ircevent"
	"log"
	"time"
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

type message struct {
	Sender string
	Text   string
	Sent   time.Time
}

var ircConn *irc.Connection
var BotAPI *tgbotapi.BotAPI

func onWelcome(event *irc.Event) {
	ircConn.Join("#tag")
	time.Sleep(100 * time.Millisecond)

	ircConn.SendRawf("%s %s :%s", "PRIVMSG", "#tag", "TorontoCryptoBot has joined the channel.")
}

func onPrivateMessage(event *irc.Event) {
	go func(event *irc.Event) {
		fmt.Print("\n\n\n\n\n")
		fmt.Print("\n\n\n\n\n")

		fmt.Print(event.Nick) // contains the sender

		fmt.Print("\n\n\n\n\n")

		fmt.Print(event.Message()) //contains the message

		fmt.Print("\n\n\n\n\n")

		fmt.Print(event.Arguments[0]) // contains the channel

		fmt.Print("\n\n\n\n\n")
		fmt.Print("\n\n\n\n\n")
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		//BotAPI.Send(msg)
	}(event)

}

func setUpIRCConnection() <-chan message {
	channel := make(chan message)

	go func() {
		ircConn = irc.IRC("Abot", "Abot")
		ircConn.VerboseCallbackHandler = true

		ircConn.Connect("irc.oftc.net:6667")

		ircConn.AddCallback("001", onWelcome)
		ircConn.AddCallback("PRIVMSG", onPrivateMessage)

		ircConn.Loop()
	}()

	//ircConn.UseTLS = ircBot.UseTLS
	//ircConn.TLSConfig = &tls.Config{
	//	ServerName: ircBot.Host,
	//}

	ircConn.VerboseCallbackHandler = ircBot.Debug
	return channel

}

func main() {
	setUpIRCConnection()

	x := 0

	for {
		x += 0
	}

	return // BEGIN TELEGRAM CODE BELOW

	TOKEN := ""

	bot, e := tgbotapi.NewBotAPI(TOKEN)

	if e != nil {
		log.Panic(e)
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

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}

}
