package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
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
	leaderboard, err := bot.GenerateLeaderboardEmbed(interaction.GuildID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bot.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "No leaderboard found",
							Description: "There is no leaderboard for this guild.",
							Color:       0xFF0000,
						},
					},
				}})
		}
		return err
	}

	return bot.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{leaderboard},
		},
	})
}
