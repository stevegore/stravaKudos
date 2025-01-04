package bot

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/stevegore/stravaKudos/parser"
)

func (s *StravaBot) GetMyProfile(c *parser.Client) (jsonData string) {

	var headers = map[string]string{}

	headers["authorization"] = "access_token " + s.authToken

	jsonData, statusCode := c.MakeRequest(s.MapUrls["my_profile"], "GET", "", headers)

	if statusCode == 401 {
		s.toAuth(c)
		jsonData = s.GetMyProfile(c)
	}

	var result map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &result)

	if err != nil {
		slog.Error("couldn't unmarshal my profile", "err", slog.String("error", err.Error()))
		return
	}

	if _, ok := result["id"]; ok {
		s.athleteId = strconv.Itoa(int(result["id"].(float64)))
	}

	return
}
