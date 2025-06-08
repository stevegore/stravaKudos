package bot

import "github.com/stevegore/stravaKudos/stravaapi"

// StravaBot represents a bot for interacting with the Strava API.
type StravaBot struct {
	authToken    string
	athleteID    string
	Friends      []string
	FriendsInfo  map[string]string
	apiEndpoints map[string]string
	Client       *stravaapi.Client
}

// NewStravaBot creates a new instance of StravaBot.
func NewStravaBot() *StravaBot {
	bot := &StravaBot{}

	var siteDomain = "https://cdn-1.strava.com"
	var langParam = "hl=en"

	bot.apiEndpoints = map[string]string{
		"auth_url":    siteDomain + "/api/v3/oauth/internal/token?" + langParam,
		"my_profile":  siteDomain + "/api/v3/athlete?" + langParam,
		"friends_url": siteDomain + "/api/v3/athletes/{ATHLETE-ID}/friends?" + langParam,
		"feed_url":    siteDomain + "/api/v3/feed/athlete/{ATHLETE-ID}",
		"feed_param":  "?photo_sizes[]=240&single_entity_supported=true&modular=true&" + langParam,
		"kudos_url":   siteDomain + "/api/v3/activities/{ACTIVITIES-ID}/kudos?" + langParam,
	}

	bot.ReadAuthToken()

	if bot.authToken != "" {
		bot.Client = stravaapi.NewClient(bot.authToken)
	} else {
		bot.Client = stravaapi.NewClient()
	}

	bot.Friends = []string{}
	bot.FriendsInfo = make(map[string]string)

	return bot
}

// Run starts the Strava bot.
// It performs authentication, retrieves the user's profile and friends, and then gives kudos to friends' activities.
func (s *StravaBot) Run() {
}
