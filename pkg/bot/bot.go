package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"

	"github.com/muncomputersciencesociety/elf/pkg/config"
)

type Bot struct {
	Session *discordgo.Session
	Config  config.Config
}

func (bot *Bot) Start() error {
	err := bot.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening bot session: %w", err)
	}

	log.Info().Msg("Elf is now running. Press Ctrl-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Info().Msg("Stopping Elf...")

	err = bot.Session.Close()
	if err != nil {
		return fmt.Errorf("error disconnecting from Discord: %w", err)
	}

	return nil
}

func New(config config.Config) (*Bot, error) {
	bot := &Bot{
		Config: config,
	}

	log.Debug().Msg("Creating Discord session")
	session, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	bot.Session = session

	return bot, nil
}
