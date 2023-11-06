package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/nint8835/elf/pkg/api"
	"github.com/nint8835/elf/pkg/bot"
	"github.com/nint8835/elf/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading config")
		os.Exit(1)
	}

	botInst, err := bot.New(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Error creating bot")
		os.Exit(1)
	}

	apiInst := api.New(botInst)

	log.Info().Msg("Starting bot")
	go func() {
		err := botInst.Start()
		if err != nil {
			log.Error().Err(err).Msg("Error starting bot")
			os.Exit(1)
		}
	}()

	log.Info().Msg("Starting API")
	err = apiInst.Start()
	if err != nil {
		log.Error().Err(err).Msg("Error starting leaderboard API")
	}

	log.Info().Msg("Stopping bot")
	botInst.Stop()
}
