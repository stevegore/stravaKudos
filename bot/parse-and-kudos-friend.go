package bot

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"math/rand"

	"github.com/stevegore/stravaKudos/parser"
)

func (s *Strava) ParseAndKudosFriend(friendId string) {

	c := &parser.Client{}

	c.InitWebClient()

	c.SetUserAgent("Strava/33.0.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.004)")

	var headers = map[string]string{}
	headers["authorization"] = "access_token " + s.authToken

	var feedFolowerUrl = strings.ReplaceAll(s.MapUrls["feed_url"], "{ATHLETE-ID}", friendId) + s.MapUrls["feed_param"]

	jsonData, statusCode := c.MakeRequest(feedFolowerUrl, "GET", "", headers)

	if statusCode != 200 {
		log.Fatalf("Status from get friend feed (friend => %s) request no HTTP_OK | statusCode => %d", friendId, statusCode)
	}

	var results []parser.FeedEntry

	err := json.Unmarshal([]byte(jsonData), &results)
	c.CheckError(err)

	c.ToLog("Getting activities for", s.FriendsInfo[friendId], "(", friendId, ")")
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	for _, result := range results {

		item := result.Item

		if item.StartDate.After(oneWeekAgo) && !item.HasKudoed {
			c.ToLog("Giving kudos to", item.Name, "by", s.FriendsInfo[friendId], "(", friendId, ")")

			s.kudosFriend(c, item.ID)

			n := rand.Intn(10) // n will be between 0 and 10
			time.Sleep(time.Duration(n) * time.Second)
		}
	}
}
