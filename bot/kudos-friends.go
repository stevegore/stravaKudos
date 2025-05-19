package bot

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/stevegore/stravaKudos/parser"
)

func (s *StravaBot) kudosActivity(activity parser.Item) {
	activityId := activity.ID
	if activityId == 0 {
		slog.Error("activity id is 0, can't give kudos")
		return
	}
	slog.Debug("giving kudos", "activityId", activityId, "activityName", activity.Name)

	var kudosUrl = strings.ReplaceAll(s.MapUrls["kudos_url"], "{ACTIVITIES-ID}", strconv.FormatInt(activityId, 10))

	resp, statusCode := s.Client.MakeRequest(kudosUrl, "POST", "", nil)

	if statusCode != 201 {
		slog.Error("couldn't give kudos", "statusCode", statusCode, "body", resp, "activityId", activityId)
		return
	}
	slog.Debug("successfully gave kudos", "activityId", activityId, "activityName", activity.Name)
}
