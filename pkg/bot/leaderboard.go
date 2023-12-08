package bot

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kanrichan/resvg-go"

	"github.com/nint8835/elf/pkg/adventofcode"
	"github.com/nint8835/elf/pkg/database"
)

//go:embed Inter.ttf
var interFontData []byte

func (bot *Bot) GenerateLeaderboardImage(guildId string) (*discordgo.File, error) {
	var guild database.Guild
	if tx := bot.Database.First(&guild, "guild_id = ?", guildId); tx.Error != nil {
		return nil, fmt.Errorf("error fetching guild details: %w", tx.Error)
	}

	if guild.LeaderboardID == nil {
		return nil, errors.New("no leaderboard id set")
	}

	worker, err := resvg.NewDefaultWorker(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error creating resvg worker: %w", err)
	}
	defer worker.Close()

	fontdb, err := worker.NewFontDBDefault()
	if err != nil {
		return nil, fmt.Errorf("error creating fontdb: %w", err)
	}
	defer fontdb.Close()

	err = fontdb.LoadFontData(interFontData)
	if err != nil {
		return nil, fmt.Errorf("error loading font data: %w", err)
	}

	err = fontdb.SetSansSerifFamily("Inter")
	if err != nil {
		return nil, fmt.Errorf("error setting sans-serif font family: %w", err)
	}

	err = fontdb.SetSerifFamily("Inter")
	if err != nil {
		return nil, fmt.Errorf("error setting serif font family: %w", err)
	}

	tree, err := worker.NewTreeFromData([]byte(`<svg version="1.1"
     width="300" height="200"
     xmlns="http://www.w3.org/2000/svg">

  <rect width="100%" height="100%" fill="red" />

  <circle cx="150" cy="100" r="80" fill="green" />

  <text x="150" y="125" font-size="60" text-anchor="middle" fill="white">SVG</text>

</svg>
`), &resvg.Options{})
	if err != nil {
		return nil, fmt.Errorf("error creating tree: %w", err)
	}
	defer tree.Close()

	tree.ConvertText(fontdb)

	treeWidth, treeHeight, err := tree.GetSize()
	if err != nil {
		return nil, fmt.Errorf("error getting tree size: %w", err)
	}

	pixmap, err := worker.NewPixmap(uint32(treeWidth), uint32(treeHeight))
	if err != nil {
		return nil, fmt.Errorf("error creating pixmap: %w", err)
	}

	err = tree.Render(resvg.TransformIdentity(), pixmap)
	if err != nil {
		return nil, fmt.Errorf("error rendering leaderboard: %w", err)
	}

	png, err := pixmap.EncodePNG()
	if err != nil {
		return nil, fmt.Errorf("error encoding PNG: %w", err)
	}

	return &discordgo.File{
		Name:   "leaderboard.png",
		Reader: bytes.NewReader(png),
	}, nil
}

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

	for i, member := range leaderboardEntries[:int(math.Min(float64(len(leaderboardEntries)), 20))] {
		stars := "`"

		for dayNumber := 1; dayNumber <= 25; dayNumber++ {
			day, ok := member.CompletionDayLevel[strconv.Itoa(dayNumber)]
			if !ok {
				stars += " "
				continue
			}
			_, star1 := day["1"]
			_, star2 := day["2"]
			if star1 && star2 {
				stars += "★"
			} else if star1 || star2 {
				stars += "☆"
			}
		}

		stars = strings.TrimRight(stars, " ")

		stars += "`"

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
