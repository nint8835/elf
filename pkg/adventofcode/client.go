package adventofcode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type cacheKey struct {
	LeaderboardID string
	Event         string
}

type CachedLeaderboard struct {
	Leaderboard Leaderboard
	RetrievedAt time.Time
}

// Client is the client for the Advent of Code API.
type Client struct {
	client *http.Client

	cache     map[cacheKey]CachedLeaderboard
	cacheLock *sync.Mutex
}

// GetLeaderboard returns the leaderboard for a given private leaderboard ID and event.
func (client *Client) GetLeaderboard(leaderboardId string, event string) (CachedLeaderboard, error) {
	client.cacheLock.Lock()
	defer client.cacheLock.Unlock()

	requestCacheKey := cacheKey{
		LeaderboardID: leaderboardId,
		Event:         event,
	}

	cachedLeaderboard, hasCachedVal := client.cache[requestCacheKey]
	if hasCachedVal && time.Since(cachedLeaderboard.RetrievedAt) < time.Minute*15 {
		log.Debug().Msg("Leaderboard requested, but cache not expired. Using cached value.")
		return cachedLeaderboard, nil
	}

	resp, err := client.client.Get(fmt.Sprintf("https://adventofcode.com/%s/leaderboard/private/view/%s.json", event, leaderboardId))
	if err != nil {
		return CachedLeaderboard{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	var leaderboard Leaderboard
	err = json.NewDecoder(resp.Body).Decode(&leaderboard)
	if err != nil {
		return CachedLeaderboard{}, fmt.Errorf("error decoding response: %w", err)
	}

	cachedLeaderboard = CachedLeaderboard{
		Leaderboard: leaderboard,
		RetrievedAt: time.Now(),
	}

	client.cache[requestCacheKey] = cachedLeaderboard

	return cachedLeaderboard, nil
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
		client:    &http.Client{Jar: jar},
		cache:     map[cacheKey]CachedLeaderboard{},
		cacheLock: &sync.Mutex{},
	}, nil
}
