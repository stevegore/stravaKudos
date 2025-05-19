package bot

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/stevegore/stravaKudos/parser"
)

func (s *StravaBot) kudosActivity(logger *slog.Logger, activity parser.Item) {
	activityId := activity.ID
	if activityId == 0 {
		logger.Error("activity id is 0, can't give kudos")
		return
	}
	logger.Info("giving kudos")

	var kudosUrl = strings.ReplaceAll(s.MapUrls["kudos_url"], "{ACTIVITIES-ID}", strconv.FormatInt(activityId, 10))

	resp, statusCode := s.Client.MakeRequest(kudosUrl, "POST", "", nil)

	if statusCode != 201 {
		logger.Error("couldn't give kudos", "statusCode", statusCode, "body", resp)
		return
	}
	logger.Debug("successfully gave kudos")
}
