package cmd

import (
	"log/slog"
	"time"

	"github.com/kohinigeee/DiscordGemuboBot/botmanager"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
	"github.com/kohinigeee/DiscordGemuboBot/slashapi"
)

func SlashApply() {

	slashapi.InitEnvs()

	commands := botmanager.InitialSlashCommands()

	for _, command := range commands {
		err := slashapi.ApplySlashCommand(&command.Command)
		if err != nil {
			mylogger.L().Error("Error applying slash command", slog.String("err", err.Error()))
		}
		time.Sleep(2 * time.Second)
	}
}
