package botmanager

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/gemubo"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

var (
	logger *slog.Logger
)

func init() {
	logger = mylogger.L()
}

type HandlerName string
type DiscorBotdHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, manager *BotManager)

type BotManager struct {
	OkReaction          string
	NoReaction          string
	batchDurationMinu   int
	Session             *discordgo.Session
	BotUserInfo         *discordgo.User
	LastNotionBatchTime time.Time
	intrHandlersMap     map[HandlerName]DiscorBotdHandler
	appHandlersMap      map[HandlerName]DiscorBotdHandler
	modalHandlerMap     map[HandlerName]DiscorBotdHandler
	templates           map[string]*gemubo.Template
	presets             map[string]*gemubo.Preset
	bosyuMsgs           map[string]*gemubo.GemuboMessage
}

func NewBotManager(s *discordgo.Session) *BotManager {
	manager := &BotManager{
		Session:           s,
		intrHandlersMap:   make(map[HandlerName]DiscorBotdHandler),
		appHandlersMap:    make(map[HandlerName]DiscorBotdHandler),
		modalHandlerMap:   make(map[HandlerName]DiscorBotdHandler),
		templates:         make(map[string]*gemubo.Template),
		presets:           make(map[string]*gemubo.Preset),
		bosyuMsgs:         make(map[string]*gemubo.GemuboMessage),
		batchDurationMinu: 3,
		OkReaction:        "👍",
		NoReaction:        "🙏",
	}

	// Add initial slash commands handler
	for _, slash := range InitialSlashCommands() {
		manager.AddAppHandler(HandlerName(slash.Command.Name), slash.Handler)
	}

	// Add initial interaction commands handler
	for _, interact := range InitialInteractCommands() {
		manager.AddIntrHandler(HandlerName(interact.Name), interact.Handler)
	}

	// Add initial modal commands handler
	for _, modal := range InitialDiscordModalCommands() {
		manager.AddModalHandler(HandlerName(modal.Name), modal.Handler)
	}

	manager.Session.AddHandler(manager.onInteractionCreate)
	return manager
}

func (manager *BotManager) cleanMmory() {
	mp := make(map[string]*gemubo.GemuboMessage, len(manager.bosyuMsgs))
	for key, value := range manager.bosyuMsgs {
		mp[key] = value
	}
	manager.bosyuMsgs = mp
}

func (manager *BotManager) notionLoop() {
	ticker := time.NewTicker(time.Minute * time.Duration(manager.batchDurationMinu))
	defer ticker.Stop()

	logger := mylogger.L()

	cleanCnt := 0

	logger.Info("Notion Loop started")
	for {
		<-ticker.C
		now := time.Now().UTC()

		nowStr := now.Add(time.Hour * 9).Format("2006-01-02 15:04:05")
		logger.Info("Notion Batch", slog.String("Date", nowStr), slog.String("CleanCnt", fmt.Sprintf("%d", cleanCnt)))

		for _, msg := range manager.bosyuMsgs {
			if msg.StartTime.Before(now) {
				manager.BosyuNotion(msg.GemuboID)
			}
		}

		manager.LastNotionBatchTime = now

		cleanCnt++
		if cleanCnt > 5 {
			manager.cleanMmory()
			cleanCnt = 0
		}
	}
}

func (manager *BotManager) Start() {
	manager.BotUserInfo = manager.Session.State.User
	go manager.notionLoop()
}

func (manager *BotManager) AddAppHandler(name HandlerName, handler DiscorBotdHandler) {
	manager.appHandlersMap[name] = handler
}

func (manager *BotManager) AddIntrHandler(name HandlerName, handler DiscorBotdHandler) {
	manager.intrHandlersMap[name] = handler
}

func (manager *BotManager) AddModalHandler(name HandlerName, handler DiscorBotdHandler) {
	manager.modalHandlerMap[name] = handler
}

func (manager *BotManager) AddTemplate(template *gemubo.Template) {
	manager.templates[template.Name] = template
}

func (manager *BotManager) AddPreset(preset *gemubo.Preset) {
	manager.presets[preset.Name] = preset
}

func (manager *BotManager) AddScheduledBosyuMsg(msg *gemubo.GemuboMessage) {
	manager.bosyuMsgs[msg.GemuboID] = msg
}

func (manager *BotManager) DeleteTemplate(name string) (exist bool, deletPresetNames []string) {
	deletPresetNames = make([]string, 0)
	_, exist = manager.templates[name]
	if !exist {
		return
	}

	for _, preset := range manager.presets {
		if preset.Template.Name == name {
			manager.DeletePreset(preset.Name)
			deletPresetNames = append(deletPresetNames, preset.Name)
		}
	}
	return
}

func (manager *BotManager) DeletePreset(name string) (exist bool) {
	_, exist = manager.presets[name]
	delete(manager.presets, name)
	return
}

func (manager *BotManager) BosyuNotion(gemuboId string) {
	gmsg, exist := manager.bosyuMsgs[gemuboId]
	if !exist {
		return
	}

	msgContent := ""
	embedTitle := "全員集合～!"

	okUsers, err := manager.Session.MessageReactions(gmsg.ChannelID, gmsg.MessageID, manager.OkReaction, 10, "", "")
	if err != nil {
		logger.Error("Error getting reactions", slog.String("err", err.Error()))
		return
	}

	msgContent += gmsg.Author.Mention() + " "
	for _, user := range okUsers {
		if user.Bot {
			continue
		}
		msgContent += user.Mention() + " "
	}
	msgContent += "\n"

	embed := &discordgo.MessageEmbed{
		Title:       embedTitle,
		Description: "",
		Color:       0x00F1AA,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: manager.BotUserInfo.AvatarURL("20"),
		},
	}

	options := &discordgo.MessageSend{
		Content: msgContent,
		Reference: &discordgo.MessageReference{
			MessageID: gmsg.MessageID,
		},
		Embed: embed,
	}

	_, err = manager.Session.ChannelMessageSendComplex(gmsg.ChannelID, options)
	if err != nil {
		logger.Error("Error sending message", slog.String("err", err.Error()))
		errmsg := fmt.Sprintf("開始通知の送信に失敗しました\n(ID:%s)", gmsg.GemuboID)
		manager.SendErrorMessage(gmsg.ChannelID, "Bosyu Notion", errmsg, nil)
		return
	}

	delete(manager.bosyuMsgs, gemuboId)
}

func (manager *BotManager) SendNormalMessage(channelID string, title string, msg string, fileds []*discordgo.MessageEmbedField) {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: msg,
		Color:       0x00ff00,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: manager.BotUserInfo.AvatarURL("20"),
		},
	}

	if fileds != nil {
		embed.Fields = fileds
	}

	_, err := manager.Session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		logger.Error("Error sending message", slog.String("err", err.Error()))
	}
}

func (manager *BotManager) SendErrorMessage(channelID string, title string, msg string, fileds []*discordgo.MessageEmbedField) {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: msg,
		Color:       0xff0000,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: manager.BotUserInfo.AvatarURL("20"),
		},
	}

	if fileds != nil {
		embed.Fields = fileds
	}

	_, err := manager.Session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		logger.Error("Error sending message", slog.String("err", err.Error()))
	}
}

// Discord event handler
func (manager *BotManager) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.Type {

	case discordgo.InteractionMessageComponent:
		handlerName := HandlerName(i.MessageComponentData().CustomID)
		handler, ok := manager.intrHandlersMap[handlerName]
		if !ok {
			logger.Warn("Handler not found", slog.String("handlerName", string(handlerName)))
			return
		}
		handler(s, i, manager)

	case discordgo.InteractionModalSubmit:
		handlerName := HandlerName(i.ModalSubmitData().CustomID)
		handler, ok := manager.modalHandlerMap[handlerName]
		if !ok {
			logger.Warn("Handler not found", slog.String("handlerName", string(handlerName)))
			return
		}
		handler(s, i, manager)

	case discordgo.InteractionApplicationCommand:
		handlerName := HandlerName(i.ApplicationCommandData().Name)
		handler, ok := manager.appHandlersMap[handlerName]
		if !ok {
			logger.Warn("Handler not found", slog.String("handlerName", string(handlerName)))
			return
		}
		handler(s, i, manager)

	default:

	}
}
