package discord

import (
	"archecord/internal/app/archeage"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Example: "!채권" or "!채권-누이"
	order := strings.Split(m.Content, "-")

	if order[0] == "!채권" {
		serverName := "모르페우스"
		if len(order) >= 2 {
			serverName = order[1]
		}
		bondInfo := archeage.BondParser(archeage.ServerNameMap[serverName])
		route := archeage.RecommendRoute(bondInfo)
		_, err := s.ChannelMessageSend(m.ChannelID, "[추천]  :  "+route)
		if err != nil {
			fmt.Println("error opening connection, ", err)
		}
	}
}