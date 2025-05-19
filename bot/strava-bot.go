package bot

type StravaBot struct {
	authToken   string
	athleteId   string
	Friends     []string
	FriendsInfo map[string]string
	MapUrls     map[string]string
}

func NewStravaBot() *StravaBot {
	s := &StravaBot{}

	var siteDomain = "https://cdn-1.strava.com"
	var langParam = "hl=en"

	s.MapUrls = map[string]string{
		"auth_url":    siteDomain + "/api/v3/oauth/internal/token?" + langParam,
		"my_profile":  siteDomain + "/api/v3/athlete?" + langParam,
		"friends_url": siteDomain + "/api/v3/athletes/{ATHLETE-ID}/friends?" + langParam,
		"feed_url":    siteDomain + "/api/v3/feed/athlete/{ATHLETE-ID}",
		"feed_param":  "?photo_sizes[]=240&single_entity_supported=true&modular=true&" + langParam,
		"kudos_url":   siteDomain + "/api/v3/activities/{ACTIVITIES-ID}/kudos?" + langParam,
	}

	s.ReadAuthToken()

	s.Friends = []string{}
	s.FriendsInfo = make(map[string]string)

	return s
}
