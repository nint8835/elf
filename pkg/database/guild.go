package database

import "gorm.io/gorm"

type Guild struct {
	gorm.Model

	// GuildID is the Discord ID for this guild.
	GuildID string `gorm:"index"`
	// LeaderboardCode is the Advent of Code leaderboard code for this guild's leaderboard.
	LeaderboardCode *string
	// LeaderboardID is the unique ID for the Advent of Code leaderboard for this guild.
	LeaderboardID *string
	// ChannelID is the ID of the channel to post leaderboard updates to.
	ChannelID *string
	// EnableAPI configures whether the data for this leaderboard should be republished via our API.
	EnableAPI bool
}
