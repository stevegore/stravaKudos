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
	friendLogger := slog.With("friend", friendName, "friendId", friendId)

	friendLogger.Debug("getting activities")
	var feedFollowerUrl = strings.ReplaceAll(s.MapUrls["feed_url"], "{ATHLETE-ID}", friendId) + s.MapUrls["feed_param"]

	jsonData, statusCode := s.Client.MakeRequest(feedFollowerUrl, "GET", "", nil)

	if statusCode != 200 {
		friendLogger.Error("couldn't get friend's activities", "statusCode", statusCode, "body", jsonData)
		return
	}

	var results []parser.FeedEntry

	err := json.Unmarshal([]byte(jsonData), &results)
	if err != nil {
		friendLogger.Error("couldn't unmarshal friend's activities", slog.Any("error", err))
		return
	}

	if len(results) == 0 {
		friendLogger.Debug("no activities found for friend")
		return
	}

	recentEnough := time.Now().AddDate(0, 0, -2)
	activitiesKudoed := 0
	for _, result := range results {

		item := result.Item
		activityLogger := friendLogger.With("activity", item.Name, "date", item.StartDate.Format("2006-01-02 15:04"), "id", item.ID)
		if item.StartDate.After(recentEnough) && !item.HasKudoed {
			s.kudosActivity(activityLogger, item)
			activitiesKudoed++

			n := 10 + rand.Intn(20)
			time.Sleep(time.Duration(n) * time.Second)
		}
	}
	if activitiesKudoed == 0 {
		friendLogger.Debug("no new activities to kudos for friend")
	}
}
