package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func leaderboardCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate, _ struct{}) error {
	interactionData, err := botInst.GenerateLeaderboardMessage(interaction.GuildID)
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
		log.Error().Err(err).Msg("Error generating leaderboard")
		return err
	}

	interactionData.Flags = discordgo.MessageFlagsEphemeral

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: interactionData,
	})
}
