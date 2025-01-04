// build for linux-amd64
// GOOS=linux GOARCH=amd64 go build strava-kudos.go

// build for mac-amd64
// GOOS=darwin GOARCH=amd64 go build strava-kudos.go

package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/stevegore/stravaKudos/bot"
	"github.com/stevegore/stravaKudos/parser"

	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {

	c := &parser.Client{}

	c.InitWebClient()

	c.SetUserAgent("Strava/33.0.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.004)")

	initDebug(c)

	s := &bot.Strava{}

	var siteDomain = "https://strava.com"
	var langParam = "hl=en"

	s.MapUrls = map[string]string{
		"auth_url":    siteDomain + "/api/v3/oauth/internal/token?" + langParam,
		"my_profile":  siteDomain + "/api/v3/athlete?" + langParam,
		"friends_url": siteDomain + "/api/v3/athletes/{ATHLETE-ID}/friends?" + langParam,
		"feed_url":    siteDomain + "/api/v3/feed/athlete/{ATHLETE-ID}",
		"feed_param":  "?photo_sizes[]=240&single_entity_supported=true&modular=true&" + langParam,
		"kudos_url":   siteDomain + "/api/v3/activities/{ACTIVITIES-ID}/kudos?" + langParam,
	}

	s.ReadAuthToken()

	for {
		s.GetMyProfile(c)
		s.GetMyFriends(c)

		c.ToLog("Friends: ", s.Friends)

		// Create a semaphore channel to limit concurrency to 8
		sem := make(chan struct{}, 8)

		for _, friendId := range s.Friends {
			sem <- struct{}{} // Acquire a slot

			go func(friendId string) {
				defer func() { <-sem }() // Release the slot

				s.ParseAndKudosFriend(friendId)
			}(friendId)
		}

		// Wait for all goroutines to finish
		for i := 0; i < cap(sem); i++ {
			sem <- struct{}{}
		}

		c.ToLog(" THE END LOOP ")
		time.Sleep(2 * time.Hour)

	}

}

func initDebug(c *parser.Client) {
	debugString, DebugEnvExists := os.LookupEnv("DEBUG")

	if DebugEnvExists {
		debug, err := strconv.ParseBool(debugString)
		if err != nil {
			log.Panicf("func initDebug(): failed convert 'debugString' to bool : %s", err)
		}
		c.SetDebug(debug)
	} else {
		c.SetDebug(true)
	}
}
