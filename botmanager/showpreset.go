package botmanager

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/lib"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

const (
	showPresetPrestNameOptionName string = "preset_name"
)

func ShowPresetHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	const HandlerTitle = "Show Preset"

	logger := mylogger.L()
	logger.Debug("ShowPresetHandler started")

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})

	if err != nil {
		mylogger.L().Warn("Error empty responding to interaction", "err", err.Error())
	}

	options := i.ApplicationCommandData().Options
	presetNameOpt := lib.GetOptionByName(options, showPresetPrestNameOptionName)

	if presetNameOpt == nil {
		content := ""
		for _, preset := range manager.presets {
			content += "-\t" + preset.Name + "\n"
		}
		fileds := make([]*discordgo.MessageEmbedField, 0)
		fileds = append(fileds, &discordgo.MessageEmbedField{
			Name:   "プリセット一覧",
			Value:  content,
			Inline: false,
		})
		manager.SendNormalMessage(i.ChannelID, HandlerTitle, "", fileds)
		logger.Info("Show Preset List", slog.String("Content", content))
		return
	}

	presetName := presetNameOpt.Value.(string)
	preset, exist := manager.presets[presetName]
	if !exist {
		errmsg := fmt.Sprintf("プリセット `%s` は存在しません", presetName)
		manager.SendErrorMessage(i.ChannelID, HandlerTitle, errmsg, nil)
		return
	}

	presetContent := ""
	keys := lib.MapKyes(preset.Params)

	for _, key := range keys {
		presetContent += fmt.Sprintf("-\t%s: %s\n", key, preset.Params[key])
	}
	fileds := make([]*discordgo.MessageEmbedField, 0, 2)
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   "プリセット名",
		Value:  preset.Name + "\n",
		Inline: true,
	})
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   "参照テンプレート名",
		Value:  preset.Template.Name + "\n",
		Inline: true,
	})
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   "プリセット内容",
		Value:  presetContent + "\n",
		Inline: true,
	})
	manager.SendNormalMessage(i.ChannelID, HandlerTitle, "", fileds)

	logger.Info("Show Preset", slog.String("presetName", presetName))
}
