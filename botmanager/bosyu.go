package botmanager

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/gemubo"
	"github.com/kohinigeee/DiscordGemuboBot/lib"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

const (
	bosyuPresetNameOptionName = "preset_name"
	BosyuModalName            = "bosyu_handler_modal"
)

type bosyuHandlerMemory struct {
	presetName string
	eventID    string
	author     *discordgo.User
}

var (
	mem = map[string]bosyuHandlerMemory{}
)

func BosyuHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {
	const handlerName = "BosyuHandler"
	logger := mylogger.L()

	logger.Debug(handlerName+"started", "ID", i.Interaction.ID)

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
		errmsg := fmt.Sprintf("プリセット `%s` は存在しません", presetName)
		manager.SendErrorMessage(i.ChannelID, handlerName, errmsg, nil)
		return
	}

	mem[i.ID] = bosyuHandlerMemory{
		presetName: presetName,
		eventID:    i.ID,
		author:     i.Member.User,
	}
	//メモリリークを防ぐために，一定時間後に削除する
	go func() {
		time.Sleep(45 * time.Second)
		delete(mem, i.ID)
	}()

	modalData := makeBosyuHandlerModal(preset, i)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: modalData,
	})

	if err != nil {
		logger.Error(handlerName+"Error responding modal to interaction", "err", err)
	}

	logger.Debug(handlerName+"ended", slog.String("prestName", presetName), slog.String("ID", i.Interaction.ID))
}

func BosyuModalHandler(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager) {

	const handlerName = "Bosyu"
	logger := mylogger.L()
	logger.Debug(handlerName+" started", slog.String("ID", i.Interaction.ID))

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})

	if err != nil {
		mylogger.L().Warn(handlerName+"Error empty responding to interaction", slog.String("err", err.Error()))
	}

	data := i.ModalSubmitData()
	eventId := lib.GetModalDataValue(&data, 0, 0)
	normalPlaceValuesStr := lib.GetModalDataValue(&data, 2, 0)
	spetialPlaceValuesStr := lib.GetModalDataValue(&data, 3, 0)

	memory, exist := mem[eventId]
	if !exist {
		logger.Error(handlerName + "memory not found")
		errmsg := fmt.Sprintf("イベントID `%s` は不正な値です", eventId)
		manager.SendErrorMessage(i.ChannelID, handlerName, errmsg, nil)
		return
	}

	preset, exist := manager.presets[memory.presetName]
	if !exist {
		logger.Error(handlerName + "preset not found")
		errmsg := "プリセット `" + memory.presetName + "` は存在しません"
		manager.SendErrorMessage(i.ChannelID, handlerName, errmsg, nil)
		return
	}

	additionalPlaceValues := lib.ParsePlaceHolders(normalPlaceValuesStr)
	spetialPlaceValues := lib.ParsePlaceHolders(spetialPlaceValuesStr)

	for key, val := range spetialPlaceValues {
		additionalPlaceValues[key] = val
	}

	if missing, ok := preset.IsFullTemplateParams(additionalPlaceValues); !ok {
		errmsg := fmt.Sprintf("プレースホルダーが不足しています.`%v`", missing)
		manager.SendErrorMessage(i.ChannelID, handlerName, errmsg, nil)
		return
	}

	bosyuMsg, err := preset.MakeMessage(additionalPlaceValues, i.ChannelID, i.GuildID, memory.author)

	if err != nil {
		logger.Error(handlerName+"Error making message", "err", err)
		errmsg := fmt.Sprintf("メッセージの作成に失敗しました.`%v`", err)
		manager.SendErrorMessage(i.ChannelID, handlerName, errmsg, nil)
		return
	}

	embed := gemubo.MakeEmbedBosyuMessage(bosyuMsg, manager.BotUserInfo)
	content := "@everyone\n"
	msgObj := &discordgo.MessageSend{
		Content: content,
		Embed:   embed,
	}

	dmsg, err := s.ChannelMessageSendComplex(i.ChannelID, msgObj)

	if err != nil {
		logger.Error(handlerName+"Error sending message", "err", err)
		errmsg := fmt.Sprintf("メッセージの送信に失敗しました.`%v`", err)
		manager.SendErrorMessage(i.ChannelID, handlerName, errmsg, nil)
		return
	}

	bosyuMsg.MessageID = dmsg.ID
	if bosyuMsg.StartTime != nil {
		manager.AddScheduledBosyuMsg(bosyuMsg)
		logger.Info(handlerName+" Scheduled message has been added", slog.String("GemuboID", bosyuMsg.GemuboID))
	}

	s.MessageReactionAdd(i.ChannelID, dmsg.ID, manager.OkReaction)
	s.MessageReactionAdd(i.ChannelID, dmsg.ID, manager.NoReaction)

	logger.Debug(handlerName+" ended", slog.String("ID", i.Interaction.ID))
}

func makeBosyuHandlerModal(preset *gemubo.Preset, i *discordgo.InteractionCreate) *discordgo.InteractionResponseData {

	components := []discordgo.MessageComponent{}

	components = append(components, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID: "bosyu_eventId",
				Label:    "イベントID(変更しないでください)",
				Value:    i.ID,
				Style:    discordgo.TextInputShort,
			},
		},
	})

	components = append(components, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID: "bosyu_template_content",
				Label:    "参照テンプレートです(変更しないでください)",
				Value:    preset.Template.Content,
				Style:    discordgo.TextInputParagraph,
			},
		},
	})

	normalHalodersContent := ""
	keys := lib.MapKyes[any](preset.Template.PlaeceHolders)

	for _, key := range keys {
		if gemubo.IsSpetialPlaceHolder(key) {
			continue
		}
		if val, exist := preset.Params[key]; exist {
			normalHalodersContent += key + "=" + val + "\n"
		} else {
			normalHalodersContent += key + "=\n"
		}
	}

	components = append(components, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID: "bosyu_normal_holders",
				Label:    "プレースホルダーの値を入力してください.",
				Value:    normalHalodersContent,
				Style:    discordgo.TextInputParagraph,
				Required: true,
			},
		},
	})

	speialHalodersContent := ""
	for _, placeHolder := range gemubo.SpetialPlaceHolders() {
		if val, exist := preset.Params[placeHolder.Name]; exist {
			speialHalodersContent += placeHolder.Name + "=" + val + "\n"
		} else {
			speialHalodersContent += placeHolder.Name + "=\n"
		}
	}

	components = append(components, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID: "bosyu_speial_holders",
				Label:    "特殊プレースホルダーの値を入力してください.時刻は`hh:mm`の形式で入力してください",
				Value:    speialHalodersContent,
				Style:    discordgo.TextInputParagraph,
				Required: true,
			},
		},
	})

	return &discordgo.InteractionResponseData{
		CustomID:   BosyuModalName,
		Title:      "プレースホルダ―入力",
		Components: components,
	}
}
