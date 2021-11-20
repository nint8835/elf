package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const InteractionResponseFlagEphemeral = 1 << 6

// Commands is a list of all Discord commands belonging to this bot.
var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "leaderboard",
		Description: "Displays the current leaderboard for this guild.",
	},
}

var commandHandlers = map[string]func(*Bot, *discordgo.InteractionCreate) error{
	"leaderboard": leaderboardCommand,
}

func leaderboardCommand(bot *Bot, interaction *discordgo.InteractionCreate) error {
	err := bot.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This command is not yet implemented",
			Flags:   InteractionResponseFlagEphemeral,
		},
	})
	if err != nil {
		return fmt.Errorf("error responding to command: %w", err)
	}

	return nil
}
