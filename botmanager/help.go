package botmanager

import (
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

func HelpHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	const HandlerName string = "Help"
	logger := mylogger.L()
	logger.Debug("HelpHandler started", "ID", i.Interaction.ID)
	defer logger.Debug("HelpHandler ended", "ID", i.Interaction.ID)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})
	if err != nil {
		logger.Warn("Error empty responding to interaction", slog.String("err", err.Error()))
	}

	const prefixTabs = 1

	commandContent := ""
	for _, command := range InitialSlashCommands() {
		commandContent += "-\t" + "```" + command.Command.Name + "```" + "\n"
		commandContent += addTabParagraph(command.Command.Description, prefixTabs)
	}

	embedFileds := []*discordgo.MessageEmbedField{}

	embedFileds = append(embedFileds, &discordgo.MessageEmbedField{
		Name:   "スラッシュコマンド一覧\t\t(`/コマンド名`で使用可能)",
		Value:  commandContent,
		Inline: false,
	})

	placeHolderStr := ""

	// $START_TIMEについての説明
	holerDescription := ""
	holerDescription += "募集の通知時刻を設定します．\n"
	holerDescription += "時刻の指定は`hh:mm`形式でしてください. \n"
	holerDescription += "(`指定なし`または`Now`で置換すると通知を行いません) \n"
	placeHolderStr += "-\t```$START_TIME```" + "\n"
	placeHolderStr += addTabParagraph(holerDescription, prefixTabs)

	// $TITLEについての説明
	holerDescription = ""
	holerDescription += "募集のタイトルを設定します．\n"
	placeHolderStr += "-\t```$TITLE```" + "\n"
	placeHolderStr += addTabParagraph(holerDescription, prefixTabs)

	// $IMAGE_URLについての説明
	holerDescription = ""
	holerDescription += "募集に表示する画像のURLを設定します." + "\n"
	placeHolderStr += "-\t```$IMAGE_URL```" + "\n"
	placeHolderStr += addTabParagraph(holerDescription, prefixTabs)

	embedFileds = append(embedFileds, &discordgo.MessageEmbedField{
		Name:   "特殊プレースホルダについて",
		Value:  placeHolderStr,
		Inline: false,
	})

	manager.SendNormalMessage(i.ChannelID, HandlerName, "", embedFileds)
}

func addTabParagraph(content string, tabs int) string {
	sentences := strings.Split(content, "\n")
	prefix := strings.Repeat("\t", tabs)
	para := ""
	for _, sensentence := range sentences {
		if sensentence == "" {
			continue
		}
		para += prefix + sensentence + "\n"
	}

	return para
}
