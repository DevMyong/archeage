package main

import (
	"archecord/internal/app/archeage"
	"archecord/internal/app/discord"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a new discord session
	dg, err := discordgo.New("Bot " + archeage.Token)
	if err != nil {
		fmt.Println("error creating Discord session,")
		log.Fatal(err)
	}

	// Add Handler
	dg.AddHandler(discord.MessageHandler)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open the Websocket
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,")
		log.Fatal(err)
	}

	log.Println("Bot is running. Press CTRL-C to exit.")

	// Create channel to keep program
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
