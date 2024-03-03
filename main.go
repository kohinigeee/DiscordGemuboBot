package main

import (
	"flag"
	"log/slog"

	"github.com/kohinigeee/DiscordGemuboBot/cmd"
	"github.com/kohinigeee/DiscordGemuboBot/mylogger"
)

//必要権限
//Send Message
//Send Messages in Threads
//Manager Messages
//Manager Threads
//Read Message History
//Mentoion Everyone

func main() {
	logger := mylogger.L()

	var mode string

	flag.StringVar(&mode, "mode", "boot", "boot or slashapply")
	flag.Parse()

	switch mode {
	case "boot":
		logger.Info("Starting bot boot mode")
		cmd.BotBoot()
	case "slashapply":
		logger.Info("Starting slash apply mode")
		cmd.SlashApply()
	default:
		logger.Error("Invalid mode name", slog.String("mode", mode))
		flag.Usage()
	}

}
