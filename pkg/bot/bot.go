package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/muncomputersciencesociety/elf/pkg/adventofcode"
	"github.com/muncomputersciencesociety/elf/pkg/config"
	"github.com/muncomputersciencesociety/elf/pkg/database"
)

type Bot struct {
	Session            *discordgo.Session
	Config             config.Config
	Database           *gorm.DB
	AdventOfCodeClient *adventofcode.Client
}

func (bot *Bot) handleCommand(interaction *discordgo.InteractionCreate) {
	commandName := interaction.ApplicationCommandData().Name
	handler, ok := commandHandlers[commandName]
	if !ok {
		log.Error().Str("command", commandName).Msg("Got interaction event for unknown command")
		return
	}
	err := handler(bot, interaction)
	if err != nil {
		log.Error().Str("command", commandName).Err(err).Msg("Error handling command")
	}
}

func (bot *Bot) onInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		bot.handleCommand(interaction)
	default:
		log.Warn().Interface("interaction", interaction).Msg("Got unknown interaction event")
	}

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
	session.AddHandler(bot.onInteractionCreate)
	bot.Session = session

	log.Debug().Msg("Creating DB instance")
	db, err := database.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("error creating DB instance: %w", err)
	}
	bot.Database = db

	log.Debug().Msg("Creating Advent of Code client")
	client, err := adventofcode.NewClient(config.AdventOfCodeSession)
	if err != nil {
		return nil, fmt.Errorf("error creating Advent of Code client: %w", err)
	}
	bot.AdventOfCodeClient = client

	return bot, nil
}
