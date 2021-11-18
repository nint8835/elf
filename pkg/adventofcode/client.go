package adventofcode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// Client is the client for the Advent of Code API.
type Client struct {
	client *http.Client
}

// GetLeaderboard returns the leaderboard for a given private leaderboard ID and event.
func (client *Client) GetLeaderboard(leaderboardId string, event string) (Leaderboard, error) {
	resp, err := client.client.Get(fmt.Sprintf("https://adventofcode.com/%s/leaderboard/private/view/%s.json", event, leaderboardId))
	if err != nil {
		return Leaderboard{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	var leaderboard Leaderboard
	err = json.NewDecoder(resp.Body).Decode(&leaderboard)
	if err != nil {
		return Leaderboard{}, fmt.Errorf("error decoding response: %w", err)
	}

	return leaderboard, nil
}

// NewClient initializes a new client.
func NewClient(session string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar: %w", err)
	}
	aocUrl, err := url.Parse("https://adventofcode.com")
	if err != nil {
		return nil, fmt.Errorf("error parsing Advent of Code URL: %w", err)
	}

	jar.SetCookies(aocUrl, []*http.Cookie{
		{
			Name:  "session",
			Value: session,
		},
	})
	return &Client{
		client: &http.Client{Jar: jar},
	}, nil
}
