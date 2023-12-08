package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func leaderboardCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate, _ struct{}) error {
	leaderboard, err := botInst.GenerateLeaderboardEmbed(interaction.GuildID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "No leaderboard found",
							Description: "There is no leaderboard for this guild.",
							Color:       0xFF0000,
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
		}
		return err
	}

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{leaderboard},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}
