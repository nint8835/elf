package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/muncomputersciencesociety/elf/pkg/bot"
	"github.com/muncomputersciencesociety/elf/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading config")
		os.Exit(1)
	}

	bot, err := bot.New(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Error creating bot")
		os.Exit(1)
	}

	err = bot.Start()
	if err != nil {
		log.Error().Err(err).Msg("Error starting bot")
		os.Exit(1)
	}
}
