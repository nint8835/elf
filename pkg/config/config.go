package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// LogLevel is the level of logging to use.
	LogLevel string `default:"info" split_words:"true"`

	// DatabasePath is the path to a SQLite database file to use for data storage.
	DatabasePath string `default:"elf.sqlite" split_words:"true"`

	// DiscordToken is the token of the Discord bot to run as.
	DiscordToken string `required:"true" split_words:"true"`
	// DiscordAppID is the ID of the app owning the bot user.
	DiscordAppID string `required:"true" split_words:"true"`
	// DiscordGuildID is the ID of the guild to register commands on, for development.
	DiscordGuildID string `split_words:"true"`
	// AdventOfCodeSession is the session cookie for Advent of Code of the bot user to use for this bot instance.
	AdventOfCodeSession string `required:"true" split_words:"true"`

	// AdventOfCodeEvent is the event to fetch details from.
	AdventOfCodeEvent string `default:"2024" split_words:"true"`
	// UpdateSchedule is a cron schedule expression (in UTC) denoting when to update leaderboards
	UpdateSchedule string `default:"30 1,15 1-25 Dec *" split_words:"true"`

	// ApiBindAddr is the address the leaderboard API should listen on.
	ApiBindAddr string `default:":10000" split_words:"true"`
}

func Load() (Config, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Info().Err(err).Msg("Error loading .env file")
	}

	var config Config
	err = envconfig.Process("elf", &config)

	if err != nil {
		return Config{}, fmt.Errorf("error parsing Elf config: %w", err)
	}

	logLevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing log level: %w", err)
	}
	zerolog.SetGlobalLevel(logLevel)

	return config, nil
}
