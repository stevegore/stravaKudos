package bot

import (
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"
)

// GetMyFriends retrieves the list of friends for the authenticated user.
// It makes a GET request to the Strava API and returns the list of friends as a string.
func (s *StravaBot) GetMyFriends() {
	var myFolowersUrl = strings.ReplaceAll(s.MapUrls["friends_url"], "{ATHLETE-ID}", s.athleteId)

	jsonData, statusCode := s.Client.MakeRequest(myFolowersUrl, "GET", "", nil)

	if statusCode != 200 {
		slog.Error("couldn't get friends", "statusCode", statusCode, "body", jsonData)
		return
	}

	var results []map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &results)
	if err != nil {
		slog.Error("couldn't unmarshal friends", slog.Any("error", err))
		return
	}

	for _, result := range results {
		if idFloat, ok := result["id"].(float64); ok {
			friendsId := strconv.Itoa(int(idFloat))

			firstname, fnOk := result["firstname"].(string)
			lastname, lnOk := result["lastname"].(string)
			username := ""
			if fnOk {
				username += firstname
			}
			if lnOk {
				if fnOk {
					username += " "
				}
				username += lastname
			}
			if username == "" {
				slog.Warn("friend found with no name", "friendId", friendsId)
				username = "Unnamed Friend"
			}

			s.Friends = append(s.Friends, friendsId)
			s.FriendsInfo[friendsId] = username
		} else {
			slog.Warn("friend data found without a valid ID")
		}
	}
}
