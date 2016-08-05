package main

import (
	//"crypto/tls"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/thoj/go-ircevent"
	"log"
	"time"
	"os"
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
		fmt.Print("\n\n")
		fmt.Print("SETTING UP THE TELEGRAM CONNECTION")

		TOKEN := ""

		if TOKEN == "" {
			fmt.Print("Telegram Token not set. Exiting.")
			os.Exit(1)
		}

		bot, e := tgbotapi.NewBotAPI(TOKEN)

		fmt.Print("Bot authenticated attempted")

		if e != nil {
			log.Panic(e)
		}

		bot.Debug = true

		log.Printf("Authorized on account %s", bot.Self.UserName)

		fmt.Print("Getting new updates from telegram")
		u := tgbotapi.NewUpdate(0)

		u.Timeout = 60

		updates, error := bot.GetUpdatesChan(u)

		if error != nil {
			log.Panic(error)
		}

		// read from channel
		go func() {
			for {
				fmt.Print("READing Telegram channel")
				message := <-writeChannel
				fmt.Print(message)
				msg := tgbotapi.NewMessage(-122277940, message)
				bot.Send(msg)

			}
		}()

		for update := range updates {
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] %s %s", update.Message.From.UserName, update.Message.Text, update.Message.Chat.ID)

			go func() {
				fmt.Print("Writing to IRC Channel")
				message := fmt.Sprintf("[%s]: %s", update.Message.From.UserName, update.Message.Text)
				ircReadChannel <- message
			}()

			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID

			//bot.Send(msg)
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
			ircConn.SendRawf("%s %s :%s", "PRIVMSG", "#tag", message)

		}
	}()

	fmt.Print("Entering the infinite loop for irc events")
	fmt.Print("\n\n\n\n")
	time.Sleep(100 * time.Millisecond)

	go func() {
		ircConn.Loop()
	}()

	//ircConn.UseTLS = ircBot.UseTLS
	//ircConn.TLSConfig = &tls.Config{
	//	ServerName: ircBot.Host,
	//}

	//ircConn.VerboseCallbackHandler = ircBot.Debug

}

func main() {
	ircWriteChannel := make(chan string)
	ircReadChannel := make(chan string)
	telegramWriteChannel := make(chan string)
	telegramReadChannel := make(chan string)

	setUpTelegramConnection(telegramReadChannel, telegramWriteChannel, ircReadChannel, ircWriteChannel)
	setUpIRCConnection(ircReadChannel, ircWriteChannel, telegramReadChannel, telegramWriteChannel)

	for {
		// keep main program running
	}

	return // BEGIN TELEGRAM CODE BELOW

}
