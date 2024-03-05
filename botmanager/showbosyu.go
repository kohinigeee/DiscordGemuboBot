package botmanager

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

func ShowBosyuHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	const handlerName = "Show Bosyu"
	logger := mylogger.L()

	logger.Debug(handlerName+" started", "ID", i.Interaction.ID)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})

	if err != nil {
		logger.Warn("Error empty responding to interaction", slog.String("err", err.Error()))
	}

	content := ""
	for _, gmsg := range manager.bosyuMsgs {
		messageLink := ""
		messageLink = "https://discord.com/channels/" + gmsg.GuildId + "/" + gmsg.ChannelID + "/" + gmsg.MessageID

		content += fmt.Sprintf("-\tID: %s ([Content](<%s>))\n", gmsg.GemuboID, messageLink)
	}

	title := "募集一覧"
	manager.SendNormalMessage(i.ChannelID, title, content, nil)

	// インタラクションを終了
	logger.Debug(handlerName+" ended", "ID", i.Interaction.ID)
}
