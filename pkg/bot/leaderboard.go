package bot

import (
	"errors"
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"

	"github.com/muncomputersciencesociety/elf/pkg/adventofcode"
	"github.com/muncomputersciencesociety/elf/pkg/database"
)

func (bot *Bot) GenerateLeaderboardEmbed(guildId string) (*discordgo.MessageEmbed, error) {
	var guild database.Guild
	if tx := bot.Database.First(&guild, "guild_id = ?", guildId); tx.Error != nil {
		return nil, fmt.Errorf("error fetching guild details: %w", tx.Error)
	}

	if guild.LeaderboardID == nil {
		return nil, errors.New("no leaderboard id set")
	}

	leaderboard, err := bot.AdventOfCodeClient.GetLeaderboard(*guild.LeaderboardID, bot.Config.AdventOfCodeEvent)
	if err != nil {
		return nil, fmt.Errorf("error fetching leaderboard: %w", err)
	}

	leaderboardEmbed := &discordgo.MessageEmbed{
		Title:  "Leaderboard",
		URL:    fmt.Sprintf("https://adventofcode.com/%s/leaderboard/private/view/%s", bot.Config.AdventOfCodeEvent, *guild.LeaderboardID),
		Color:  0x007152,
		Fields: []*discordgo.MessageEmbedField{},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Join code: %s", *guild.LeaderboardCode),
		},
	}

	leaderboardEntries := []adventofcode.LeaderboardMember{}

	for _, member := range leaderboard.Members {
		leaderboardEntries = append(leaderboardEntries, member)
	}

	sort.Slice(leaderboardEntries, func(i, j int) bool {
		return leaderboardEntries[i].LocalScore > leaderboardEntries[j].LocalScore
	})

	for i, member := range leaderboardEntries[:10] {
		leaderboardEmbed.Fields = append(leaderboardEmbed.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%d. %s", i+1, member.Name),
			Value: fmt.Sprintf("%d points", member.LocalScore),
		})
	}

	return leaderboardEmbed, nil
}
