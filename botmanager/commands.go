package botmanager

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/slashapi"
)

type SlashCommand struct {
	Command slashapi.SlashCommandJson
	Handler DiscorBotdHandler
}

type InteractCommand struct {
	Name    string
	Handler DiscorBotdHandler
}

type DiscordModalCommand struct {
	Name    string
	Handler DiscorBotdHandler
}

var (
	SlashCommands = []SlashCommand{
		{
			Command: slashapi.SlashCommandJson{
				Name:        "ping",
				Description: "ping pong",
			},
			Handler: pingPong,
		},
		{
			Command: slashapi.SlashCommandJson{
				Name:        "gemubo_set_template",
				Description: "テンプレートをセットします",
				Options:     nil,
			},
			Handler: SetTemplateHandler,
		},
		{
			Command: slashapi.SlashCommandJson{
				Name:        "gemubo_show_templates",
				Description: "登録されているテンプレートを表示します(テンプレート名を指定すると，そのテンプレートの詳細を表示します)",
				Options: []slashapi.SlashCommandOptionJson{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Description: "指定なしは，テンプレート一覧を表示します",
						Name:        templateNameOptionName,
						Required:    false,
					},
				},
			},
			Handler: ShowTemplateHandler,
		},
		{
			Command: slashapi.SlashCommandJson{
				Name:        "gemubo_set_preset",
				Description: "プリセットをセットします",
				Options: []slashapi.SlashCommandOptionJson{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Description: "参照するテンプレート名を指定してください(一覧は`/gemubo_show_templates`で確認できます)",
						Name:        setPresetTemplateOptionName,
						Required:    true,
					},
				},
			},
			Handler: SetPresetHandler,
		},
		{
			Command: slashapi.SlashCommandJson{
				Name:        "gemubo_show_presets",
				Description: "登録されているプリセットを表示します(プリセット名を指定すると，そのプリセットの詳細を表示します)",
				Options: []slashapi.SlashCommandOptionJson{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Description: "指定なしは，プリセット一覧を表示します",
						Name:        showPresetPrestNameOptionName,
						Required:    false,
					},
				},
			},
			Handler: ShowPresetHandler,
		},
	}

	//------------------------------------------
	InteractCommands = []InteractCommand{}

	//------------------------------------------
	DiscordModalCommands = []DiscordModalCommand{
		{
			Name:    SetTemplateModalName,
			Handler: SetTemplateModalHandler,
		},
		{
			Name:    SetPresetModalName,
			Handler: SetPrestModalHandler,
		},
	}
)

func InitialSlashCommands() []SlashCommand {
	return SlashCommands
}

func InitialInteractCommands() []InteractCommand {
	return InteractCommands
}

func InitialDiscordModalCommands() []DiscordModalCommand {
	return DiscordModalCommands
}
