package cmd

import (
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

		guildId, _ := cmd.Flags().GetString("guild-id")
		leaderboardCode, _ := cmd.Flags().GetString("leaderboard-code")
		leaderboardId, _ := cmd.Flags().GetString("leaderboard-id")
		channelId, _ := cmd.Flags().GetString("channel-id")

		guild := &database.Guild{
			GuildID:         guildId,
			LeaderboardCode: &leaderboardCode,
			LeaderboardID:   &leaderboardId,
			ChannelID:       &channelId,
		}

		tx := bot.Database.Create(guild)
		checkError(tx.Error, "Error creating test guild")

		log.Info().Msg("Test guild added")
	},
}

func init() {
	rootCmd.AddCommand(testdataCmd)

	testdataCmd.Flags().String("guild-id", "", "guild ID to add to the database")
	testdataCmd.Flags().String("leaderboard-code", "", "leaderboard code to add to the database")
	testdataCmd.Flags().String("leaderboard-id", "", "leaderboard ID to add to the database")
	testdataCmd.Flags().String("channel-id", "", "channel ID to add to the database")
}
