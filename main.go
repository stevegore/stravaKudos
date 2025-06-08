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
		slog.Error("debug env variable must be a boolean, will assume true", slog.Any("error", err))
		debug = true
	}
	if debug {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format(time.RFC3339))
			}
			return a
		},
	}))
	slog.SetDefault(logger)

	slog.Info("starting Strava kudos bot")
	slog.Debug("debug mode enabled")
}

// main is the entry point of the application.
// It initializes and runs the Strava bot.
func main() {

	stravaBot := bot.NewStravaBot()
	cronJob := cron.New()

	performKudosCycle(stravaBot) // Do it at start up

	cronSchedules := []string{
		"30 9 * * 0",   // Sunday at 9:30
		"44 8 * * 1-5", // Weekdays at 8:44
		"24 10 * * 6",  // Saturday at 10:24
		"4 12 * * *",   // Every day at 12:04
	}

	for _, schedule := range cronSchedules {
		_, err := cronJob.AddFunc(schedule, func() {
			delayRandomly()
			performKudosCycle(stravaBot)
		})
		if err != nil {
			slog.Error("failed to schedule task", slog.Any("error", err))
			os.Exit(1)
		}
	}

	cronJob.Start()

	// Keep the program running
	select {}
}

func delayRandomly() {
	// Delay a random number of minutes between 1-10 minutes
	delay := time.Duration(rand.Intn(10)+1) * time.Minute
	slog.Debug("pausing for random interval", "duration", delay.String())
	time.Sleep(delay)
}

func performKudosCycle(stravaBot *bot.StravaBot) {
	stravaBot.GetMyProfile()
	stravaBot.GetMyFriends()
	if len(stravaBot.Friends) == 0 {
		slog.Info("no friends found")
		return
	}
	slog.Debug("retrieved friends", "friend count", strconv.Itoa(len(stravaBot.Friends)))

	// Create a semaphore channel to limit concurrency to 8
	sem := make(chan struct{}, 8)

	for _, friendId := range stravaBot.Friends {
		sem <- struct{}{} // Acquire a slot

		go func(id string) {
			defer func() { <-sem }() // Release the slot

			stravaBot.ParseAndKudosFriend(id)
		}(friendId)
	}

	// Wait for all goroutines to finish
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
}
