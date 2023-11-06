package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "elf",
	Short: "An Advent of Code <-> Discord integration bot.",
}

func Execute() {
	err := rootCmd.Execute()
	checkError(err, "Failed to run")
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}
