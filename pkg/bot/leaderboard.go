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
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/bwmarrin/discordgo"
	"github.com/kanrichan/resvg-go"

	"github.com/nint8835/elf/pkg/adventofcode"
	"github.com/nint8835/elf/pkg/database"
)

var noStarsStyle = "stroke:rgb(100,100,100);stroke-width:1;fill:none;"
var oneStarStyle = "fill:rgb(150,150,150);"
var twoStarsStyle = "fill:rgb(241,150,0);"

//go:embed JetBrainsMono.ttf
var jetBrainsMonoFontData []byte

//go:embed leaderboard.svg
var leaderboardTemplateSvg string

var leaderboardTemplate = template.Must(template.New("leaderboard").Funcs(sprig.FuncMap()).Parse(leaderboardTemplateSvg))

type LeaderboardEntry struct {
	Username string
	Days     []string
}

type LeaderboardTemplateData struct {
	Entries []LeaderboardEntry
	Event   string
}

func templateLeaderboardSvg(leaderboard adventofcode.CachedLeaderboard) ([]byte, error) {
	currentDay := time.Now().UTC().Day()

	data := LeaderboardTemplateData{
		Event: leaderboard.Leaderboard.Event,
	}

	leaderboardEntries := make([]adventofcode.LeaderboardMember, 0)

	for _, member := range leaderboard.Leaderboard.Members {
		leaderboardEntries = append(leaderboardEntries, member)
	}

	sort.Slice(leaderboardEntries, func(i, j int) bool {
		return leaderboardEntries[i].LocalScore > leaderboardEntries[j].LocalScore
	})

	for _, member := range leaderboardEntries[:int(math.Min(float64(len(leaderboard.Leaderboard.Members)), 10))] {
		var days []string

		for dayNumber := 1; dayNumber <= currentDay; dayNumber++ {
			day, ok := member.CompletionDayLevel[strconv.Itoa(dayNumber)]
			if !ok {
				days = append(days, noStarsStyle)
				continue
			}
			_, star1 := day["1"]
			_, star2 := day["2"]
			if star1 && star2 {
				days = append(days, twoStarsStyle)
			} else if star1 || star2 {
				days = append(days, oneStarStyle)
			} else {
				days = append(days, noStarsStyle)
			}
		}

		data.Entries = append(data.Entries, LeaderboardEntry{
			Username: member.Name,
			Days:     days,
		})
	}

	var buffer bytes.Buffer
	err := leaderboardTemplate.Execute(&buffer, data)
	if err != nil {
		return nil, fmt.Errorf("error executing leaderboard template: %w", err)
	}

	return buffer.Bytes(), nil
}

func renderSvg(content []byte) ([]byte, error) {
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

	err = fontdb.LoadFontData(jetBrainsMonoFontData)
	if err != nil {
		return nil, fmt.Errorf("error loading font data: %w", err)
	}

	err = fontdb.SetMonospaceFamily("JetBrains Mono")
	if err != nil {
		return nil, fmt.Errorf("error setting monospace font family: %w", err)
	}

	err = fontdb.SetSerifFamily("Inter")
	if err != nil {
		return nil, fmt.Errorf("error setting serif font family: %w", err)
	}

	tree, err := worker.NewTreeFromData(content, &resvg.Options{})
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

	return png, nil
}

func (bot *Bot) GenerateLeaderboardMessage(guildId string) (*discordgo.InteractionResponseData, error) {
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

	templatedLeaderboard, err := templateLeaderboardSvg(leaderboard)
	if err != nil {
		return nil, fmt.Errorf("error templating leaderboard: %w", err)
	}

	png, err := renderSvg(templatedLeaderboard)
	if err != nil {
		return nil, fmt.Errorf("error rendering leaderboard: %w", err)
	}

	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Leaderboard",
				URL:   fmt.Sprintf("https://adventofcode.com/%s/leaderboard/private/view/%s", bot.Config.AdventOfCodeEvent, *guild.LeaderboardID),
				Color: 0x007152,
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://leaderboard.png",
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Join code: %s", *guild.LeaderboardCode),
				},
				Timestamp: leaderboard.RetrievedAt.Format(time.RFC3339),
			},
		},
		Files: []*discordgo.File{
			{
				Name:   "leaderboard.png",
				Reader: bytes.NewReader(png),
			},
		},
	}, nil
}
