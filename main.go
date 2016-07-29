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
	ircConn       *irc.Connection
	Nickname      string       // The nick the bot will use
	Username      string       // The IRC username the bot will use
	Password      string       // Server password
	Host          string       // The IRC Server hostname
	Port          int          // The IRC Server port
	UseTLS        bool         // Should we connect using TLS?
	TLSServerName string       // Must supply if above is true
	Debug         bool         // Log all IRC comms to std out
	Channel       chan message // Communication back to main thread
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
func setUpTelegramConnection(
	readChannel chan string,
	writeChannel chan string,
	ircReadChannel chan string,
	ircWriteChannel chan string) {

	go func() {
		fmt.Print("\n\n\n")
		fmt.Print("SETTING UP THE TELEGRAM CONNECTION")
		fmt.Print("\n\n\n")

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

		// read from channel
		go func() {
			for {
				fmt.Print("Writing to Telegram channel")
				message := <-writeChannel
				fmt.Print(message)
				msg := tgbotapi.NewMessage(1, message)
				bot.Send(msg)

			}
		}()

		for update := range updates {
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			go func() {
				fmt.Print("Writing to IRC Channel")
				ircWriteChannel <- update.Message.Text
			}()

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}

	}()

}

func setUpIRCConnection(
	readChannel chan string,
	writeChannel chan string,
	telegramReadChannel chan string,
	telegramWriteChannel chan string) {

	fmt.Print("\n\n\n\n")
	fmt.Print("setting up the connection to irc")

	ircConn = irc.IRC("Abbbot", "Abbbot")
	ircConn.VerboseCallbackHandler = true

	ircConn.Connect("irc.oftc.net:6667")
	//ircConn.Connect("chat.freenode.net:6667")

	fmt.Print("setting up the callbacks")
	ircConn.AddCallback("001", func(e *irc.Event) {

		ircConn.Join("#tag")
		time.Sleep(100 * time.Millisecond)

		ircConn.SendRawf("%s %s :%s", "PRIVMSG", "#tag", "TorontoCryptoBot has joined the channel.")
		return
	})

	ircConn.AddCallback("PRIVMSG", func(event *irc.Event) {
		fmt.Print(event.Nick)         // contains the sender
		fmt.Print(event.Message())    //contains the message
		fmt.Print(event.Arguments[0]) // contains the channel

		go func() {
			fmt.Print("\n\n\n")
			fmt.Print("WRITING TO TELEGRAM CHANNEL")
			telegramWriteChannel <- event.Message()
		}()

	})

	// read from channel
	go func() {
		for {
			fmt.Print("Reading from IRC channel")
			message := <-readChannel
			fmt.Print(message)

		}
	}()

	fmt.Print("Entering the infinite loop for irc events")
	fmt.Print("\n\n\n\n")

	ircConn.Loop()

	//ircConn.UseTLS = ircBot.UseTLS
	//ircConn.TLSConfig = &tls.Config{
	//	ServerName: ircBot.Host,
	//}

	//ircConn.VerboseCallbackHandler = ircBot.Debug

}

func main() {
	//	telegramChannel := setUpTelegram()  // TODO: create this channel, pass it to IRC connection
	ircWriteChannel := make(chan string)
	ircReadChannel := make(chan string)
	telegramWriteChannel := make(chan string)
	telegramReadChannel := make(chan string)

	go func() {
		setUpIRCConnection(ircReadChannel, ircWriteChannel, telegramReadChannel, telegramWriteChannel)
	}()

	go func() {
		setUpTelegramConnection(telegramReadChannel, telegramWriteChannel, ircReadChannel, ircWriteChannel)
	}()

	for {
		// keep main program running
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
