package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

type telegramBot struct {
	Token string
}

func setUpTelegramConnection(readChannel chan string, writeChannel chan string) {

	go func() {
		fmt.Print("\n\n")
		fmt.Print("SETTING UP THE TELEGRAM CONNECTION")

		TOKEN := "220802676:AAESVYAwKZipv-D7B8i8LHsAtO-fcn2b2Q4"

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

		fmt.Println("Getting new updates from telegram")
		u := tgbotapi.NewUpdate(0)

		u.Timeout = 60

		updates, error := bot.GetUpdatesChan(u)

		if error != nil {
			log.Panic(error)
		}

		// read from channel
		go func() {
			for {
				fmt.Println("READing Telegram channel")
				message := <-readChannel
				fmt.Println("MEssage was written")
				fmt.Println(message)
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
				fmt.Println("Writing to IRC Channel")
				message := fmt.Sprintf("[%s]: %s", update.Message.From.UserName, update.Message.Text)
				writeChannel <- message
			}()

			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID

			//bot.Send(msg)
		}

	}()

}
