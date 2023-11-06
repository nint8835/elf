package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/nint8835/elf/pkg/bot"
	"github.com/nint8835/elf/pkg/config"
)

var guildIdFlag = flag.String("guild-id", "", "Guild ID to register commands on. Should be used for testing.")
var globalFlag = flag.Bool("global", false, "Register commands globally.")
var unregisterFlag = flag.Bool("unregister", false, "Unregister registered commands.")

func main() {
	flag.Parse()

	if *guildIdFlag == "" && !*globalFlag {
		log.Error().Msg("--guild-id or --global must be specified.")
		os.Exit(1)
	} else if *guildIdFlag != "" && *globalFlag {
		log.Error().Msg("--guild-id and --global cannot be specified together.")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading config")
		os.Exit(1)
	}

	botInstance, err := bot.New(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Error creating bot")
		os.Exit(1)
	}

	if *unregisterFlag {
		commands, err := botInstance.Session.ApplicationCommands(cfg.DiscordAppID, *guildIdFlag)
		if err != nil {
			log.Error().Err(err).Msg("Error getting registered commands")
			os.Exit(1)
		}

		for _, command := range commands {
			err = botInstance.Session.ApplicationCommandDelete(cfg.DiscordAppID, *guildIdFlag, command.ID)
			if err != nil {
				log.Error().Interface("command", command).Err(err).Msg("Error unregistering command")
			} else {
				log.Info().Interface("command", command).Msg("Unregistered command")
			}
		}
	} else {
		for _, command := range bot.Commands {
			cmdObj, err := botInstance.Session.ApplicationCommandCreate(cfg.DiscordAppID, *guildIdFlag, command)

			if err != nil {
				log.Error().Interface("command", command).Err(err).Msg("Error registering command")
			} else {
				log.Info().Interface("command", command).Interface("cmdObj", cmdObj).Msg("Registered command")
			}
		}
	}

}
