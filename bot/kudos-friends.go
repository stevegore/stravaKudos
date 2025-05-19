package bot

import (
	"log/slog"
	"strconv"

	"strings"

	"github.com/stevegore/stravaKudos/parser"
)

func (s *StravaBot) kudosActivity(c *parser.Client, activityId int64) {
	var headers = map[string]string{}
	headers["authorization"] = "access_token " + s.authToken

	var myFolowersUrl = strings.ReplaceAll(s.MapUrls["kudos_url"], "{ACTIVITIES-ID}", strconv.Itoa(int(activityId)))

	resp, statusCode := c.MakeRequest(myFolowersUrl, "POST", "", headers)

	if statusCode != 201 {
		slog.Error("couldn't give kudos", "statusCode", statusCode, "body", resp)
		return
	}
}
