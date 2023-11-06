package bot

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/nint8835/elf/pkg/adventofcode"
	"github.com/nint8835/elf/pkg/database"
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
		Timestamp: leaderboard.RetrievedAt.Format(time.RFC3339),
	}

	leaderboardEntries := []adventofcode.LeaderboardMember{}

	for _, member := range leaderboard.Leaderboard.Members {
		leaderboardEntries = append(leaderboardEntries, member)
	}

	sort.Slice(leaderboardEntries, func(i, j int) bool {
		return leaderboardEntries[i].LocalScore > leaderboardEntries[j].LocalScore
	})

	for i, member := range leaderboardEntries[:20] {
		stars := ""

		for dayNumber := 1; dayNumber <= 25; dayNumber++ {
			day, ok := member.CompletionDayLevel[strconv.Itoa(dayNumber)]
			if !ok {
				stars += "⬛"
				continue
			}
			_, star1 := day["1"]
			_, star2 := day["2"]
			if star1 && star2 {
				stars += "🟨"
			} else if star1 || star2 {
				stars += "⬜"
			}
		}

		stars = strings.TrimRight(stars, "⬛")

		username := member.Name
		if username == "" {
			username = fmt.Sprintf("(anonymous user #%s)", member.ID)
		}

		leaderboardEmbed.Fields = append(leaderboardEmbed.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%d. %s", i+1, username),
			Value: fmt.Sprintf("%d points\n%s", member.LocalScore, stars),
		})
	}

	return leaderboardEmbed, nil
}
