// build for linux-amd64
// GOOS=linux GOARCH=amd64 go build strava-kudos.go

// build for mac-amd64
// GOOS=darwin GOARCH=amd64 go build strava-kudos.go

package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/stevegore/stravaKudos/bot"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func init() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file found")
		os.Exit(1)
	}

	level := slog.LevelInfo
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		slog.Error("debug env variable must be a boolean, will assume true", "err", slog.String("error", err.Error()))
		debug = true
	}
	if debug {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	slog.Info("starting Strava kudos bot")
	slog.Debug("debug mode enabled")
}

// main is the entry point of the application.
// It initializes and runs the Strava bot.
func main() {

	s := bot.NewStravaBot()
	cron := cron.New()

	kudosAllFriendsActivities(s) // Do it at start up

	cronSchedules := []string{
		"30 9 * * 0",   // Sunday at 9:30
		"44 8 * * 1-5", // Weekdays at 8:44
		"24 10 * * 6",  // Saturday at 10:24
		"4 12 * * *",   // Every day at 12:04
	}

	for _, schedule := range cronSchedules {
		_, err := cron.AddFunc(schedule, func() {
			delayRandomly()
			kudosAllFriendsActivities(s)
		})
		if err != nil {
			slog.Error("failed to schedule task", "error", err)
			os.Exit(1)
		}
	}

	cron.Start()

	// Keep the program running
	select {}
}

func delayRandomly() {
	// Delay a random number of minutes between 1-10 minutes
	delay := time.Duration(rand.Intn(10)+1) * time.Minute
	slog.Debug("pausing for random interval", "duration", delay.String())
	time.Sleep(delay)
}

func kudosAllFriendsActivities(s *bot.StravaBot) {
	s.GetMyProfile()
	s.GetMyFriends()
	if len(s.Friends) == 0 {
		slog.Info("no friends found")
		return
	}
	slog.Debug("retrieved friends", "friend count", strconv.Itoa(len(s.Friends)))

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
}
