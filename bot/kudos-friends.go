package bot

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/stevegore/stravaKudos/stravaapi"
)

func (s *StravaBot) kudosActivity(logger *slog.Logger, activity stravaapi.Item) {
	activityId := activity.ID
	if activityId == 0 {
		logger.Error("activity id is 0, can't give kudos")
		return
	}
	logger.Info("giving kudos")

	kudosURL := strings.ReplaceAll(s.apiEndpoints["kudos_url"], "{ACTIVITIES-ID}", strconv.FormatInt(activityId, 10))

	responseBody, statusCode := s.Client.MakeRequest(kudosURL, "POST", "", nil)

	if statusCode != 201 {
		logger.Error("couldn't give kudos", "statusCode", statusCode, "body", responseBody)
		return
	}
	logger.Debug("successfully gave kudos")
}
