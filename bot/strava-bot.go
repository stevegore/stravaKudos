package bot

type Strava struct {
	authToken   string
	athleteId   string
	Friends     []string
	FriendsInfo map[string]string
	MapUrls     map[string]string
}
