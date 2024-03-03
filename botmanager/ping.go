package botmanager

import (
	"github.com/bwmarrin/discordgo"
)

func pingPong(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
		},
	})
}
