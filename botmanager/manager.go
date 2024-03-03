package botmanager

import (
	"log/slog"

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
	Session         *discordgo.Session
	BotUserInfo     *discordgo.User
	intrHandlersMap map[HandlerName]DiscorBotdHandler
	appHandlersMap  map[HandlerName]DiscorBotdHandler
	modalHandlerMap map[HandlerName]DiscorBotdHandler
	templates       map[string]*gemubo.Template
	presets         map[string]*gemubo.Preset
}

func NewBotManager(s *discordgo.Session) *BotManager {
	manager := &BotManager{
		Session:         s,
		intrHandlersMap: make(map[HandlerName]DiscorBotdHandler),
		appHandlersMap:  make(map[HandlerName]DiscorBotdHandler),
		modalHandlerMap: make(map[HandlerName]DiscorBotdHandler),
		templates:       make(map[string]*gemubo.Template),
		presets:         make(map[string]*gemubo.Preset),
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

func (manager *BotManager) Start() {
	manager.BotUserInfo = manager.Session.State.User
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
