package bot

import (
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"math/rand"

	"github.com/stevegore/stravaKudos/parser"
)

func (s *StravaBot) ParseAndKudosFriend(friendId string) {

	c := parser.NewClient()

	friendName := s.FriendsInfo[friendId]

	var headers = map[string]string{}
	headers["authorization"] = "access_token " + s.authToken

	slog.Debug("getting activities", "friend", friendName)
	var feedFollowerUrl = strings.ReplaceAll(s.MapUrls["feed_url"], "{ATHLETE-ID}", friendId) + s.MapUrls["feed_param"]

	jsonData, statusCode := c.MakeRequest(feedFollowerUrl, "GET", "", headers)

	if statusCode != 200 {
		slog.Error("couldn't get friend's activities", "statusCode", statusCode)
		return
	}

	var results []parser.FeedEntry

	err := json.Unmarshal([]byte(jsonData), &results)
	if err != nil {
		slog.Error("couldn't unmarshal friend's activities", "err", slog.String("error", err.Error()))
		return
	}

	recentEnough := time.Now().AddDate(0, 0, -2)
	for _, result := range results {

		item := result.Item
		if item.StartDate.After(recentEnough) && !item.HasKudoed {
			slog.Info("giving kudos", "friend", friendName, "activity", item.Name, "date", item.StartDate.Format("2006-01-02 15:04"))
			s.kudosActivity(c, item.ID)

			n := 10 + rand.Intn(20) // n will be between 10 and 30
			time.Sleep(time.Duration(n) * time.Second)
		}
	}
}
