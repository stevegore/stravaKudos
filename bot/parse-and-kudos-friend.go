package bot

import (
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"math/rand"

	"github.com/stevegore/stravaKudos/stravaapi"
)

// ParseAndKudosFriend parses a friend's activity feed and gives kudos to their activities.
// It takes the friend's ID and their activity feed as input.
func (s *StravaBot) ParseAndKudosFriend(friendId string) {

	friendName := s.FriendsInfo[friendId]
	friendLogger := slog.With("friend", friendName, "friendId", friendId)

	friendLogger.Debug("getting activities")
	feedFollowerURL := strings.ReplaceAll(s.apiEndpoints["feed_url"], "{ATHLETE-ID}", friendId) + s.apiEndpoints["feed_param"]

	jsonData, statusCode := s.Client.MakeRequest(feedFollowerURL, "GET", "", nil)

	if statusCode != 200 {
		friendLogger.Error("couldn't get friend's activities", "statusCode", statusCode, "body", jsonData)
		return
	}

	var results []stravaapi.FeedEntry

	err := json.Unmarshal([]byte(jsonData), &results)
	if err != nil {
		friendLogger.Error("couldn't unmarshal friend's activities", slog.Any("error", err))
		return
	}

	if len(results) == 0 {
		friendLogger.Debug("no activities found for friend")
		return
	}

	activityCutoffDate := time.Now().AddDate(0, 0, -2)
	activitiesKudoed := 0
	for _, result := range results {

		item := result.Item
		activityLogger := friendLogger.With("activity", item.Name, "date", item.StartDate.Format("2006-01-02 15:04"), "id", item.ID)
		if item.StartDate.After(activityCutoffDate) && !item.HasKudoed {
			s.kudosActivity(activityLogger, item)
			activitiesKudoed++

			sleepDurationSeconds := 10 + rand.Intn(20)
			time.Sleep(time.Duration(sleepDurationSeconds) * time.Second)
		}
	}
	if activitiesKudoed == 0 {
		friendLogger.Debug("no new activities to kudos for friend")
	}
}
