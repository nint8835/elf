package api

import (
	"github.com/gin-gonic/gin"

	"github.com/muncomputersciencesociety/elf/pkg/bot"
)

type API struct {
	bot    *bot.Bot
	engine *gin.Engine
}

func (a *API) Start() error {
	return a.engine.Run(a.bot.Config.ApiBindAddr)
}

func New(bot *bot.Bot) *API {
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &API{
		engine: engine,
		bot:    bot,
	}
}
