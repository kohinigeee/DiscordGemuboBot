package botmanager

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/gemubo"
	"github.com/kohinigeee/DiscordGemuboBot/lib"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

const (
	bosyuPresetNameOptionName = "preset_name"
)

func BosyuHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	const handlerName = "BosyuHandler"
	logger := mylogger.L()

	options := i.ApplicationCommandData().Options
	opt := lib.GetOptionByName(options, bosyuPresetNameOptionName)

	if opt == nil {
		logger.Error(handlerName + "option not found")
		lib.SendEmptyInteractionResponse(s, i)
		return
	}

	presetName := opt.Value.(string)
	preset, exist := manager.presets[presetName]
	if !exist {
		logger.Error(handlerName + " preset not found")
		lib.SendEmptyInteractionResponse(s, i)
		return
	}

}

func makeBosyuHandlerModal(preset *gemubo.Preset) *discordgo.InteractionResponseData {

	components := []*discordgo.MessageComponent{}

	components = append(components, &discordgo.MessageActionRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID: "bosyu_preset_name",
				Label:
			}
			}

}
