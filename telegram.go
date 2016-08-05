package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type telegramBot struct {
	Token        string
	Bot          *tgbotapi.BotAPI
	Debug        bool
	ReadChannel  chan string
	WriteChannel chan string
}

func newTelegramBot(
	Token string,
	Debug bool,
	ReadChannel chan string,
	WriteChannel chan string) *telegramBot {

	return &telegramBot{
		Token:        Token,
		Debug:        Debug,
		ReadChannel:  ReadChannel,
		WriteChannel: WriteChannel,
	}
}

func (self *telegramBot) initConnection() {
	bot, e := tgbotapi.NewBotAPI(self.Token)

	if e != nil {
		log.Panic(e)
	}

	self.Bot = bot
	self.Bot.Debug = self.Debug

	fmt.Printf("Authorized on account %s", bot.Self.UserName)

	self.initReadLoop()

}

func (self *telegramBot) initReadLoop() {
	go func() {
		for {
			fmt.Println("Reading Telegram channel")
			message := <-self.ReadChannel
			fmt.Println("IRC Message was read. Sending")
			fmt.Println(message)
			// TODO: don't hardcode group ID
			msg := tgbotapi.NewMessage(-122277940, message)
			self.Bot.Send(msg)

		}
	}()
}

func (self *telegramBot) beginLoop() {
	go func() {
		fmt.Println("Getting new updates from telegram")
		u := tgbotapi.NewUpdate(0)

		u.Timeout = 60

		updates, error := self.Bot.GetUpdatesChan(u)

		if error != nil {
			log.Panic(error)
		}

		for update := range updates {
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] %s %s", update.Message.From.UserName, update.Message.Text, update.Message.Chat.ID)

			go func() {
				fmt.Println("Writing to IRC Channel")
				message := fmt.Sprintf("[%s]: %s", update.Message.From.UserName, update.Message.Text)
				self.WriteChannel <- message
			}()

			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID

			//bot.Send(msg)
		}
	}()

}
