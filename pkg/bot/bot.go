package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
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
	Scheduler          *gocron.Scheduler
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

func (bot *Bot) updateLeaderboards() {
	var guilds []database.Guild
	if tx := bot.Database.Find(&guilds); tx.Error != nil {
		log.Error().Err(tx.Error).Msg("Error getting guilds from database")
		return
	}

	for _, guild := range guilds {
		log.Info().Str("guild", guild.GuildID).Msg("Updating leaderboard")

		if guild.ChannelID == nil {
			log.Info().Str("guild", guild.GuildID).Msg("No channel set for guild, skipping")
			continue
		}

		leaderboard, err := bot.GenerateLeaderboardEmbed(guild.GuildID)
		if err != nil {
			log.Error().Err(err).Str("guild", guild.GuildID).Msg("Error generating leaderboard")
			continue
		}

		if guild.MessageID != nil {
			_, err = bot.Session.ChannelMessageEditEmbed(*guild.ChannelID, *guild.MessageID, leaderboard)
			if err != nil {
				log.Error().Err(err).Str("guild", guild.GuildID).Msg("Error editing leaderboard")
				continue
			}
		} else {
			message, err := bot.Session.ChannelMessageSendEmbed(*guild.ChannelID, leaderboard)
			if err != nil {
				log.Error().Err(err).Str("guild", guild.GuildID).Msg("Error sending leaderboard")
				continue
			}

			guild.MessageID = &message.ID
			if tx := bot.Database.Save(&guild); tx.Error != nil {
				log.Error().Err(tx.Error).Str("guild", guild.GuildID).Msg("Error saving guild to database")
			}
		}
	}
}

func (bot *Bot) Start() error {
	bot.Scheduler.StartAsync()
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

	log.Debug().Msg("Creating scheduler")
	bot.Scheduler = gocron.NewScheduler(time.UTC)
	bot.Scheduler.Cron("30 15,23 * * *").Do(bot.updateLeaderboards)

	return bot, nil
}
