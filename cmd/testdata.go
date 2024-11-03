package cmd

import (
	"github.com/charmbracelet/huh"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nint8835/elf/pkg/bot"
	"github.com/nint8835/elf/pkg/config"
	"github.com/nint8835/elf/pkg/database"
)

var testdataCmd = &cobra.Command{
	Use:   "testdata",
	Short: "Insert test data into the database",
	Args:  cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		checkError(err, "Error loading config")

		bot, err := bot.New(cfg)
		checkError(err, "Error creating bot")

		var guildId string
		var channelId string
		var leaderboardId string
		var leaderboardCode string
		var enableAPI bool

		err = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Guild ID").
					Value(&guildId),
				huh.NewInput().
					Title("Channel ID").
					Value(&channelId),
			).Title("Discord"),
			huh.NewGroup(
				huh.NewInput().
					Title("Leaderboard ID").
					Value(&leaderboardId),
				huh.NewInput().
					Title("Leaderboard Code").
					Value(&leaderboardCode),
				huh.NewConfirm().
					Title("Enable API").
					Value(&enableAPI),
			).Title("Advent of Code"),
		).WithTheme(huh.ThemeCatppuccin()).Run()
		checkError(err, "Error getting test data")

		guild := &database.Guild{
			GuildID:         guildId,
			LeaderboardCode: &leaderboardCode,
			LeaderboardID:   &leaderboardId,
			ChannelID:       &channelId,
			EnableAPI:       enableAPI,
		}

		tx := bot.Database.Create(guild)
		checkError(tx.Error, "Error creating test guild")

		log.Info().Msg("Test guild added")
	},
}

func init() {
	rootCmd.AddCommand(testdataCmd)
}
