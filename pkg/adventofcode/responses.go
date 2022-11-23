package adventofcode

import "encoding/json"

type StarDetails struct {
	GetStarTS json.Number `json:"get_star_ts"`
}

type DayCompletionStats map[string]StarDetails

type LeaderboardMember struct {
	LastStarTS         json.Number                   `json:"last_star_ts"`
	Stars              int                           `json:"stars"`
	LocalScore         int                           `json:"local_score"`
	ID                 json.Number                   `json:"id"`
	GlobalScore        int                           `json:"global_score"`
	Name               string                        `json:"name"`
	CompletionDayLevel map[string]DayCompletionStats `json:"completion_day_level"`
}

type Leaderboard struct {
	OwnerID json.Number                  `json:"owner_id"`
	Members map[string]LeaderboardMember `json:"members"`
	Event   string                       `json:"event"`
}
