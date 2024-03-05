package botmanager

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/lib"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

const (
	deletePresetOptionName = "preset_name"
)

func DeletePresetHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	const handlerName = "Delete Preset"
	logger := mylogger.L()

	logger.Debug(handlerName+" started", "ID", i.Interaction.ID)
	defer logger.Debug(handlerName+" ended", "ID", i.Interaction.ID)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})

	if err != nil {
		logger.Warn("Error empty responding to interaction", slog.String("err", err.Error()))
	}

	options := i.ApplicationCommandData().Options
	presetName := lib.GetOptionByName(options, deletePresetOptionName).StringValue()

	exist := manager.DeletePreset(presetName)
	if !exist {
		errmsg := fmt.Sprintf("プリセット `%s` は存在しません", presetName)
		manager.SendErrorMessage(i.ChannelID, handlerName, errmsg, nil)
		return
	}

	manager.SendNormalMessage(i.ChannelID, handlerName, fmt.Sprintf("プリセット `%s` を削除しました", presetName), nil)
}
