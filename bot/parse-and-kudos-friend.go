package bot

import (
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"math/rand"

	"github.com/stevegore/stravaKudos/parser"
)

// ParseAndKudosFriend parses a friend's activity feed and gives kudos to their activities.
// It takes the friend's ID and their activity feed as input.
func (s *StravaBot) ParseAndKudosFriend(friendId string) {

	friendName := s.FriendsInfo[friendId]

	slog.Debug("getting activities", "friend", friendName, "friendId", friendId)
	var feedFollowerUrl = strings.ReplaceAll(s.MapUrls["feed_url"], "{ATHLETE-ID}", friendId) + s.MapUrls["feed_param"]

	jsonData, statusCode := s.Client.MakeRequest(feedFollowerUrl, "GET", "", nil)

	if statusCode != 200 {
		slog.Error("couldn't get friend's activities", "friend", friendName, "friendId", friendId, "statusCode", statusCode, "body", jsonData)
		return
	}

	var results []parser.FeedEntry

	err := json.Unmarshal([]byte(jsonData), &results)
	if err != nil {
		slog.Error("couldn't unmarshal friend's activities", "friend", friendName, "friendId", friendId, "err", slog.String("error", err.Error()))
		return
	}

	if len(results) == 0 {
		slog.Debug("no activities found for friend", "friend", friendName, "friendId", friendId)
		return
	}

	recentEnough := time.Now().AddDate(0, 0, -2)
	activitiesKudoed := 0
	for _, result := range results {

		item := result.Item
		if item.StartDate.After(recentEnough) && !item.HasKudoed {
			slog.Info("giving kudos", "friend", friendName, "activity", item.Name, "date", item.StartDate.Format("2006-01-02 15:04"))
			s.kudosActivity(item)
			activitiesKudoed++

			n := 10 + rand.Intn(20)
			time.Sleep(time.Duration(n) * time.Second)
		}
	}
	if activitiesKudoed == 0 {
		slog.Debug("no new activities to kudos for friend", "friend", friendName, "friendId", friendId)
	}
}
