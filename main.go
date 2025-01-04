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
	"time"

	"github.com/stevegore/stravaKudos/bot"
	"github.com/stevegore/stravaKudos/parser"

	"github.com/joho/godotenv"
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

	for {
		s.GetMyProfile(c)
		s.GetMyFriends(c)
		if len(s.Friends) == 0 {
			slog.Info("no friends found")
			os.Exit(2)
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

		time.Sleep(2 * time.Hour)
	}
}
