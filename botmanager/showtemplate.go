package botmanager

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/lib"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

const (
	templateNameOptionName string = "template_name"
)

func ShowTemplateHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	const HandlerTitle = "Show Template"

	logger := mylogger.L()
	logger.Debug("ShowTemplateHandler started")

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})

	if err != nil {
		mylogger.L().Warn("Error empty responding to interaction", "err", err.Error())
	}

	options := i.ApplicationCommandData().Options
	tempNameOpt := lib.GetOptionByName(options, templateNameOptionName)

	if tempNameOpt == nil {
		content := ""
		for _, template := range manager.templates {
			content += fmt.Sprintf("-\t%s\n", template.Name)
		}
		fileds := make([]*discordgo.MessageEmbedField, 0)
		fileds = append(fileds, &discordgo.MessageEmbedField{
			Name:   "テンプレート一覧",
			Value:  content,
			Inline: false,
		})
		manager.SendNormalMessage(i.ChannelID, HandlerTitle, "", fileds)
		return
	}

	templateName := tempNameOpt.Value.(string)
	template, exist := manager.templates[templateName]
	if !exist {
		errmsg := fmt.Sprintf("テンプレート `%s` は存在しません", templateName)
		manager.SendErrorMessage(i.ChannelID, HandlerTitle, errmsg, nil)
		return
	}

	fileds := make([]*discordgo.MessageEmbedField, 0, 2)
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   "テンプレート名",
		Value:  template.Name + "\n",
		Inline: true,
	})
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   "テンプレート内容",
		Value:  template.Content + "\n",
		Inline: true,
	})
	manager.SendNormalMessage(i.ChannelID, HandlerTitle, "", fileds)
}
