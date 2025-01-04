package bot

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/stevegore/stravaKudos/parser"
)

func (s *Strava) GetMyFriends(c *parser.Client) {
	var headers = map[string]string{}

	headers["authorization"] = "access_token " + s.authToken

	var myFolowersUrl = strings.ReplaceAll(s.MapUrls["friends_url"], "{ATHLETE-ID}", s.athleteId)

	jsonData, statusCode := c.MakeRequest(myFolowersUrl, "GET", "", headers)

	if statusCode != 200 {
		log.Fatalf("Status from getMyFolowers request no HTTP_OK | statusCode => %d", statusCode)
	}

	var results []map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &results)
	c.CheckError(err)

	s.Friends = []string{}
	s.FriendsInfo = make(map[string]string)

	for _, result := range results {
		if _, ok := result["id"]; ok {

			friendsId := strconv.Itoa(int(result["id"].(float64)))

			username := result["firstname"].(string) + " " + result["lastname"].(string)

			s.Friends = append(s.Friends, friendsId)

			s.FriendsInfo[friendsId] = username
		}
	}
}
