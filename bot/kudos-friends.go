package bot

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/stevegore/stravaKudos/stravaapi"
)

// isKudosWorthy determines if an activity meets minimum requirements to be worthy of kudos.
// Returns false for swims less than 300m, runs less than 3km, and walks less than 30 minutes.
// Returns true for all other activities.
func isKudosWorthy(activity stravaapi.Item) bool {
	distance := activity.Distance
	sportType := strings.ToLower(activity.SportType)

	// Swim less than 300m
	if strings.Contains(sportType, "swim") && distance < 300 {
		return false
	}

	// Run less than 3km (3000m)
	if strings.Contains(sportType, "run") && distance < 3000 {
		return false
	}

	// Walk less than 30 mins (1800 seconds)
	if strings.Contains(sportType, "walk") && activity.ElapsedTime < 1800 {
		return false
	}

	// If the activity is less than 5 minutes, it's not ever worthy of kudos
	if activity.ElapsedTime < 300 {
		return false
	}

	// All other activities are kudos-worthy
	return true
}

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
