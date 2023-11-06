package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/nint8835/elf/pkg/bot"
	"github.com/nint8835/elf/pkg/config"
	"github.com/nint8835/elf/pkg/database"
)

var guildIdFlag = flag.String("guild-id", "", "Guild ID to add to the database")
var leaderboardCodeFlag = flag.String("leaderboard-code", "", "Leaderboard code to add to the database")
var leaderboardIdFlag = flag.String("leaderboard-id", "", "Leaderboard ID to add to the database")
var channelIdFlag = flag.String("channel-id", "", "Channel ID to add to the database")

func main() {
	flag.Parse()

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

	guild := &database.Guild{
		GuildID:         *guildIdFlag,
		LeaderboardCode: leaderboardCodeFlag,
		LeaderboardID:   leaderboardIdFlag,
		ChannelID:       channelIdFlag,
	}

	tx := bot.Database.Create(guild)

	if tx.Error != nil {
		log.Error().Err(err).Msg("Error adding test guild")
		os.Exit(1)
	}

	log.Info().Msg("Test guild added")
}
