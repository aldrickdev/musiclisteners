package main

import (
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/aldricdev/musiclisteners/internals/db"
	"github.com/aldricdev/musiclisteners/internals/types"
)

func mainUserLoop(wg *sync.WaitGroup, user types.User) {
	databasePassword := os.Getenv("APP_USER_POSTGRES_PASSWORD")
	if databasePassword == "" {
		slog.Error("Missing Environment Variable: APP_USER_POSTGRES_PASSWORD")
		wg.Done()
		return
	}

	slog.Info("Beginning Main Loop for User", "user", user.ID)

	for {
		randomSleep := rand.Intn(10) + 0
		slog.Info("User sleeping", "user_id", user.ID, "sleep", randomSleep)

		dbInstance, err := db.NewDB(databasePassword)
		if err != nil {
			slog.Error("Failed to connect to database", "error", err)
			continue
		}

		time.Sleep(time.Duration(randomSleep) * time.Second)

		song, err := dbInstance.SelectRandomSong()
		if err != nil {
			slog.Error("Failed to get a random song", "error", err.Error())
			continue
		}

		if err = dbInstance.InsertCurrentlyPlayingSongForUserTrans(user, song); err != nil {
			slog.Error("Failed to insert current playing song for user", "user_id", user.ID, "song_id", song.ID, "error", err.Error())
			continue
		}

		song, err = dbInstance.SelectCurrentlyPlayingSongForUser(user)
		if err != nil {
			slog.Error("Failed to get the current song for user", "user_id", user.ID, "error", err)
			continue
		}

		slog.Debug("Song currently playing", "user_id", user.ID, "song_id", song.ID)
		dbInstance.Connection.Close()
	}

	// Not running wg.Done() due to this being a forever loop
	// wg.Done()
}

func main() {
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	log :=	slog.New(slog.NewJSONHandler(os.Stdout, logOpts))
	slog.SetDefault(log)

	databasePassword := os.Getenv("APP_USER_POSTGRES_PASSWORD")
	if databasePassword == "" {
		slog.Error("Missing Environment Variable: APP_USER_POSTGRES_PASSWORD")
	}

	dbInstance, err := db.NewDB(databasePassword)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}
	defer dbInstance.Connection.Close()

	if err := dbInstance.Connection.Ping(); err != nil {
		slog.Error("Failed to Ping the database", "error", err)
	}
	log.Info("Connected to Database")

	users, err := dbInstance.GetAllUsers()
	if err != nil {
		slog.Error("Failed to get all users", "error", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(users))
	for _, user := range users {
		go mainUserLoop(&wg, user)
	}
	wg.Wait()
}
