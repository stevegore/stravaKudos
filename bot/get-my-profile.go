package bot

import (
	"encoding/json"
	"log/slog"
	"strconv"
)

// GetMyProfile retrieves the profile information for the authenticated user.
// It makes a GET request to the Strava API and returns the profile information as a string.
func (s *StravaBot) GetMyProfile() (jsonData string) {

	jsonData, statusCode := s.Client.MakeRequest(s.MapUrls["my_profile"], "GET", "", nil)

	if statusCode == 401 {
		s.toAuth()
		jsonData = s.GetMyProfile()
		return
	}

	if statusCode != 200 {
		slog.Error("couldn't get my profile", "statusCode", statusCode, "body", jsonData)
		return
	}

	var result map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &result)

	if err != nil {
		slog.Error("couldn't unmarshal my profile", slog.Any("error", err))
		return
	}

	if idFloat, ok := result["id"].(float64); ok {
		s.athleteId = strconv.Itoa(int(idFloat))
	} else {
		slog.Error("athlete id not found or not a number in profile response")
	}

	return
}
