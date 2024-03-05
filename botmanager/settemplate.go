package botmanager

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/gemubo"
	"github.com/kohinigeee/DiscordGemuboBot/lib"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
	"github.com/kohinigeee/mylog/clog"
)

const (
	SetTemplateModalName string = "set_template_modal"
)

func SetTemplateHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {

	logger := mylogger.L()
	logger.Debug("SetTemplateHandler started", slog.String("user", fmt.Sprintf("%+v", i.Member.User.Username)), slog.String("ID", i.ID), slog.String("InteractionID", i.Interaction.ID))

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: SetTemplateModalName,
			Title:    "Setting Template",
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "set_template_name_input",
							Label:       "テンプレート名を入力してください",
							Placeholder: "テンプレート名",
							Style:       discordgo.TextInputShort,
							Required:    true,
						},
					},
				},
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "set_template_input",
							Label:       "テンプレート内容(`$`で始めるとワードは変数になります. 例: `$GAME`)",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "[テンプレート例]\nゲーム: $GAMES\n人数 : $PEOPLE\n開始: $START_TIME\n",
							Required:    true,
						},
					},
				},
			},
		},
	})

	if err != nil {
		mylogger.L().Error("Error responding to interaction", slog.String("err", err.Error()))
	}

	logger.Debug("SetTemplateHandler ended")
}

func SetTemplateModalHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {

	logger := mylogger.L()
	logger.Debug("SetTemplateModalHandler started", slog.String("ID", i.ID), slog.String("InteractionID", i.Interaction.ID))

	data := i.ModalSubmitData()
	templateName := lib.GetModalDataValue(&data, 0, 0)
	content := lib.GetModalDataValue(&data, 1, 0)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource})

	if err != nil {
		mylogger.L().Warn("Error empty responding to interaction", slog.String("err", err.Error()))
	}

	template := gemubo.NewTemplate(templateName, content)
	manager.AddTemplate(template)

	mylogger.L().Info("Template has been set", slog.String(clog.OrderString("templateName", clog.OrderLevel1()), templateName), slog.String("params", fmt.Sprintf("%+v", template.PlaeceHolders)))

	msg := fmt.Sprintf("\nテンプレート `%s` を設定しました", templateName)
	manager.SendNormalMessage(i.ChannelID, "Set Template", msg, nil)

	logger.Debug("SetTemplateModalHandler ended")
}
