package bot

import (
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"

	"github.com/stevegore/stravaKudos/parser"
)

func (s *StravaBot) GetMyFriends(c *parser.Client) {
	var headers = map[string]string{}

	headers["authorization"] = "access_token " + s.authToken

	var myFolowersUrl = strings.ReplaceAll(s.MapUrls["friends_url"], "{ATHLETE-ID}", s.athleteId)

	jsonData, statusCode := c.MakeRequest(myFolowersUrl, "GET", "", headers)

	if statusCode != 200 {
		slog.Error("couldn't get friends", "statusCode", statusCode)
		return
	}

	var results []map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &results)
	if err != nil {
		slog.Error("couldn't unmarshal friends", "err", slog.String("error", err.Error()))
		return
	}

	for _, result := range results {
		if _, ok := result["id"]; ok {

			friendsId := strconv.Itoa(int(result["id"].(float64)))

			username := result["firstname"].(string) + " " + result["lastname"].(string)

			s.Friends = append(s.Friends, friendsId)

			s.FriendsInfo[friendsId] = username
		}
	}
}
