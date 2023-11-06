package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nint8835/elf/pkg/api"
	"github.com/nint8835/elf/pkg/bot"
	"github.com/nint8835/elf/pkg/config"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the bot",
	Args:  cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		checkError(err, "Error loading config")

		botInst, err := bot.New(cfg)
		checkError(err, "Error creating bot")

		apiInst := api.New(botInst)

		log.Info().Msg("Starting bot")
		go func() {
			err := botInst.Start()
			checkError(err, "Error starting bot")
		}()

		log.Info().Msg("Starting API")
		err = apiInst.Start()
		checkError(err, "Error starting API")

		log.Info().Msg("Stopping bot")
		botInst.Stop()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
