package main

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"time"
)

type ircBot struct {
	IRCConn       *irc.Connection
	Nickname      string      // The nick the bot will use
	Username      string      // The IRC username the bot will use
	Password      string      // Server password
	Host          string      // The IRC Server hostname
	Port          int         // The IRC Server port
	IRCChannel    string      // The IRC Channel to join
	UseTLS        bool        // Should we connect using TLS?
	TLSServerName string      // Must supply if above is true
	Debug         bool        // Log all IRC comms to std out
	ReadChannel   chan string // Communication back to main thread
	WriteChannel  chan string // Communication back to main thread
}

func newIRCBot(
	Nickname string,
	Username string,
	Host string,
	Port int,
	IRCChannel string,
	Debug bool,
	ReadChannel chan string,
	WriteChannel chan string) *ircBot {

	return &ircBot{
		Nickname:     Nickname,
		Username:     Username,
		Host:         Host,
		Port:         Port,
		IRCChannel:   IRCChannel,
		Debug:        Debug,
		ReadChannel:  ReadChannel,
		WriteChannel: WriteChannel,
	}
}

func (self *ircBot) initConnection() {
	fmt.Println("Setting up the IRC Connection")
	self.IRCConn = irc.IRC(self.Nickname, self.Username)
	self.IRCConn.VerboseCallbackHandler = self.Debug
	self.IRCConn.Connect(fmt.Sprintf("%s:%d", self.Host, self.Port))
	self.initCallbacks()

	// read from channel
	go func() {
		for {

			fmt.Println("Reading from IRC channel")
			message := <-self.ReadChannel
			fmt.Println(message)
			self.IRCConn.SendRawf("%s %s :%s", "PRIVMSG", "#tag", message)
			fmt.Println("READ FROM THE CHANNEL")

		}
	}()

}

func (self *ircBot) initCallbacks() {
	fmt.Println("Setting up the callbacks")
	// Joining channel
	self.IRCConn.AddCallback("001", func(e *irc.Event) {

		self.IRCConn.Join(self.IRCChannel)
		time.Sleep(3000 * time.Millisecond)

		self.IRCConn.SendRawf("%s %s :%s", "PRIVMSG", "#tag", "TorontoCryptoBot has joined the channel.")
		return
	})

	// Private message
	self.IRCConn.AddCallback("PRIVMSG", func(event *irc.Event) {
		fmt.Print(event.Nick)         // contains the sender
		fmt.Print(event.Message())    //contains the message
		fmt.Print(event.Arguments[0]) // contains the channel

		go func() {
			fmt.Print("\n\n\n")
			fmt.Print("WRITING TO TELEGRAM CHANNEL")
			self.WriteChannel <- event.Message()
		}()

	})
	fmt.Println("Done setting up callbacks.")

}

func (self *ircBot) beginLoop() {
	time.Sleep(15000 * time.Millisecond)
	fmt.Println("Beginning the loop")

	go func() {
		self.IRCConn.Loop()
	}()
}
