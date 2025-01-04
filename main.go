// build for linux-amd64
// GOOS=linux GOARCH=amd64 go build strava-kudos.go

// build for mac-amd64
// GOOS=darwin GOARCH=amd64 go build strava-kudos.go

package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/stevegore/stravaKudos/bot"
	"github.com/stevegore/stravaKudos/parser"

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

	slog.Info("starting strava kudos bot")
	slog.Debug("debug mode enabled")
}

func main() {

	c := parser.NewClient()
	s := bot.NewStravaBot()
	cron := cron.New()

	kudosAllFriendsActivities(c, s) // Do it at start up

	// Schedule the kudosAllFriendsActivities function to run at 8:44 AM and 12:04 every day
	_, err := cron.AddFunc("44 8 * * *", func() {
		kudosAllFriendsActivities(c, s)
	})
	if err != nil {
		slog.Error("failed to schedule task", "error", err)
		os.Exit(1)
	}

	_, err = cron.AddFunc("4 12 * * *", func() {
		kudosAllFriendsActivities(c, s)
	})
	if err != nil {
		slog.Error("failed to schedule task", "error", err)
		os.Exit(1)
	}
	cron.Start()

	// Keep the program running
	select {}
}

func kudosAllFriendsActivities(c *parser.Client, s *bot.StravaBot) {
	s.GetMyProfile(c)
	s.GetMyFriends(c)
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
