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
	LogLevel string `default:"info" split_words:"true"`

	DiscordToken        string `required:"true" split_words:"true"`
	AdventOfCodeSession string `required:"true" split_words:"true"`
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
