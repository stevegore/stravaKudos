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

	var headers = map[string]string{}
	headers["authorization"] = "access_token " + s.authToken

	var feedFolowerUrl = strings.ReplaceAll(s.MapUrls["feed_url"], "{ATHLETE-ID}", friendId) + s.MapUrls["feed_param"]

	jsonData, statusCode := c.MakeRequest(feedFolowerUrl, "GET", "", headers)

	if statusCode != 200 {
		log.Fatalf("Status from get friend feed (friend => %s) request no HTTP_OK | statusCode => %d", friendId, statusCode)
	}

	var results []map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &results)
	c.CheckError(err)

	c.ToLog("Getting activities for", s.FriendsInfo[friendId], "(", friendId, ")")

	for _, result := range results {

		item := result["item"].(map[string]interface{})

		if _, ok:= item["has_kudoed"]; ok {

			var hasKudos = item["has_kudoed"].(bool)

			if !hasKudos {
				var activitiesId = strconv.Itoa(int(result["entity_id"].(float64)))

				c.ToLog( "	new activity => ", activitiesId)

				s.kudosFollower( c, activitiesId )

				rand.Seed(time.Now().UnixNano())
				n := rand.Intn(10) // n will be between 0 and 10
				time.Sleep(time.Duration(n) * time.Second)
			}
		}
	}
}