package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muncomputersciencesociety/elf/pkg/bot"
	"github.com/muncomputersciencesociety/elf/pkg/database"
)

type API struct {
	bot    *bot.Bot
	engine *gin.Engine
}

func (a *API) Start() error {
	return a.engine.Run(a.bot.Config.ApiBindAddr)
}

func (a *API) HandleGetLeaderboard(c *gin.Context) {
	guildId := c.Param("guild_id")
	event := c.Param("event")

	var guild database.Guild
	if tx := a.bot.Database.First(&guild, "guild_id = ?", guildId); tx.Error != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Errorf("error getting guild: %w", tx.Error).Error(),
			},
		)
		return
	}

	if !guild.EnableAPI {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "API access is disabled for this guild."},
		)
		return
	}

	leaderboard, err := a.bot.AdventOfCodeClient.GetLeaderboard(*guild.LeaderboardID, event)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Errorf("error getting leaderboard: %w", err).Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"updated_at":  leaderboard.RetrievedAt,
			"leaderboard": leaderboard.Leaderboard,
		},
	)
}

func New(bot *bot.Bot) *API {
	engine := gin.New()
	engine.Use(gin.Recovery())

	apiInst := &API{
		engine: engine,
		bot:    bot,
	}

	engine.GET("/event/:event/guild/:guild_id", apiInst.HandleGetLeaderboard)

	return apiInst
}
