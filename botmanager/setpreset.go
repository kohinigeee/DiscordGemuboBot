package botmanager

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/gemubo"
	"github.com/kohinigeee/DiscordGemuboBot/lib"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

const (
	SetPresetModalName          string = "set_preset_modal"
	setPresetTemplateOptionName string = "template_name"
)

func SetPresetHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	const HandlerName string = "Set Preset"
	logger := mylogger.L()
	logger.Debug("SetPresetHandler started")

	options := i.ApplicationCommandData().Options
	templateName := options[0].StringValue()

	template, exist := manager.templates[templateName]
	if !exist {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		errmsg := fmt.Sprintf("テンプレート `%s` は存在しません", templateName)
		manager.SendErrorMessage(i.ChannelID, HandlerName, errmsg, nil)
		return
	}

	modalData := makeSetPresetModal(template)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: modalData,
	})

	if err != nil {
		logger.Error("Error responding modal to interaction", "err", err)
	}

	logger.Debug("SetPresetHandler ended", slog.String("ID", i.Interaction.ID))
}

func SetPrestModalHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	const HandlerName string = "Set Preset"
	logger := mylogger.L()
	logger.Debug("SetPrestModalHandler started")

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})

	if err != nil {
		mylogger.L().Warn("Error empty responding to interaction", slog.String("err", err.Error()))
	}

	data := i.ModalSubmitData()
	templateName := lib.GetModalDataValue(&data, 0, 0)
	presetName := lib.GetModalDataValue(&data, 1, 0)
	normalPlaceHoldersStr := lib.GetModalDataValue(&data, 2, 0)
	spetialPlaceHoldersStr := lib.GetModalDataValue(&data, 3, 0)

	placeHoldersMap := lib.ParsePlaceHolders(normalPlaceHoldersStr)
	spetialPlaceHoldersMap := lib.ParsePlaceHolders(spetialPlaceHoldersStr)

	for k, v := range spetialPlaceHoldersMap {
		placeHoldersMap[k] = v
	}

	logger.Debug("Get Modal Data", slog.String("templateName", templateName), slog.String("presetName", presetName), slog.String("placeHolders", fmt.Sprintf("%+v", placeHoldersMap)))

	templte, exist := manager.templates[templateName]
	if !exist {
		errmsg := fmt.Sprintf("テンプレート `%s` は存在しません", templateName)
		manager.SendErrorMessage(i.ChannelID, HandlerName, errmsg, nil)
		return
	}

	preset := gemubo.NewPreset(presetName, templte, placeHoldersMap)
	manager.AddPreset(preset)

	manager.SendNormalMessage(i.ChannelID, HandlerName, fmt.Sprintf("プリセット `%s` を設定しました", presetName), nil)
	logger.Info("Preset has been set", slog.String("presetName", presetName), slog.String("params", fmt.Sprintf("%+v", preset.Params)))

	logger.Debug("SetPrestModalHandler ended", slog.String("ID", i.Interaction.ID))
}

func makeSetPresetModal(template *gemubo.Template) *discordgo.InteractionResponseData {
	components := []discordgo.MessageComponent{}

	keys := lib.MapKyes[any](template.PlaeceHolders)

	components = append(components, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID: "preset_template_input",
				Label:    "テンプレート名(変更しないでください)",
				Value:    template.Name,
				Style:    discordgo.TextInputShort,
				Required: true,
			},
		},
	})

	components = append(components, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "preset_name_input",
				Label:       "プリセット名を入力してください",
				Placeholder: "プリセット名",
				Style:       discordgo.TextInputShort,
				Required:    true,
			},
		},
	})

	normalHoldersContent := ""
	for _, key := range keys {
		if gemubo.IsSpetialPlaceHolder(key) {
			continue
		}
		normalHoldersContent += fmt.Sprintf("%s=\n", key)
	}

	components = append(components, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID: "prest_normal_holders_input",
				Label:    "プレースホルダーの値を入力してください.\n指定しない場合は空欄にしてください",
				Value:    normalHoldersContent,
				Style:    discordgo.TextInputParagraph,
				Required: true,
			},
		},
	})

	speialHoldersContent := ""
	for _, placeHolder := range gemubo.SpetialPlaceHolders() {
		speialHoldersContent += fmt.Sprintf("%s=\n", placeHolder.Name)
	}

	components = append(components, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID: "prest_speial_holders_input",
				Label:    "特殊プレースホルダーの値を入力してください指定しない場合は空欄にしてください",
				Value:    speialHoldersContent,
				Style:    discordgo.TextInputParagraph,
				Required: true,
			},
		},
	})

	return &discordgo.InteractionResponseData{
		CustomID:   SetPresetModalName,
		Content:    fmt.Sprintf("Template:%s", template.Name),
		Title:      "Setting Preset",
		Components: components,
	}
}
