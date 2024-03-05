package botmanager

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/lib"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

const (
	DeleteTemplateOptionName = "template_name"
)

func DeleteTemplateHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	const handlerName = "Delete Template"

	logger := mylogger.L()
	logger.Debug(handlerName+" started", "ID", i.Interaction.ID)
	defer logger.Debug(handlerName+" ended", "ID", i.Interaction.ID)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})

	if err != nil {
		mylogger.L().Warn("Error empty responding to interaction", slog.String("err", err.Error()))
	}

	options := i.ApplicationCommandData().Options
	templateName := lib.GetOptionByName(options, DeleteTemplateOptionName).StringValue()

	exist, deltePresetNames := manager.DeleteTemplate(templateName)
	if !exist {
		errmsg := fmt.Sprintf("テンプレート `%s` は存在しません", templateName)
		manager.SendErrorMessage(i.ChannelID, handlerName, errmsg, nil)
		return
	}

	content := fmt.Sprintf("テンプレート `%s` を削除しました\n", templateName)
	deltePresets := ""
	for _, name := range deltePresetNames {
		deltePresets += fmt.Sprintf("-\t%s\n", name)
	}

	if deltePresets != "" {
		content += fmt.Sprintf("(%sを利用していた以下のプリセットを削除しました)\n", templateName)
		content += deltePresets
	}

	manager.SendNormalMessage(i.ChannelID, handlerName, content, nil)
}
