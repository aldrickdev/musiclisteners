package main

import (
	"embed"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aldricdev/musiclisteners/internals/db"
	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/aldricdev/musiclisteners/internals/utils"
)

var (
	//go:embed embed/migrations/*.sql
	migrations embed.FS

	//go:embed embed/data/*
	data embed.FS
)

const (
	// In seconds
	StartUpDelay = 10
)

func mainUserLoop2(wg *sync.WaitGroup, user types.User) {
	databasePassword := os.Getenv("APP_USER_POSTGRES_PASSWORD")
	if databasePassword == "" {
		slog.Error("Missing Environment Variable: APP_USER_POSTGRES_PASSWORD")
		wg.Done()
		return
	}
	
	dbInstance, err := db.NewDB(databasePassword)
	slog.Debug("Trying to connect to the database")
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
	}

	resultChan := make(chan db.SelectUsersResult)
	q := db.NewSelectUsers(resultChan)
	dbInstance.QueryBuffer <- q

	slog.Debug("Waiting for results")
	r := <-resultChan
	slog.Debug("Got Users", "users", r.Users)

	wg.Done()
}

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
			time.Sleep(1 * time.Minute)
			continue
		}

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

		slog.Info("User listening to new song", "user_id", user.ID, "song_id", song.ID)
		dbInstance.Connection.Close()
	}

	// Not running wg.Done() due to this being a forever loop
	// wg.Done()
}

func seedDatabase(dbInstance *db.DB) error {
	songs := utils.ImportCSVSongFromEmbededFS("embed/data/spotify_data.csv", data, utils.ExtractSongsFromCSVReader)
	slog.Debug("Extracted songs", "songs", songs)

	users := utils.GenerateUsers(5)
	slog.Debug("generated users", "users", users)

	if err := utils.SeedDB(dbInstance, songs, users); err != nil {
		slog.Error("Failed to seed database", "error", err)
		return err
	}

	return nil
}

func main() {
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	logLevel := os.Getenv("APP_LOG_LEVEL")
	switch ll := strings.ToLower(logLevel); ll {
	case "debug":
		logOpts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}

	case "error":
		logOpts = &slog.HandlerOptions{
			Level: slog.LevelError,
		}
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, logOpts))
	slog.SetDefault(log)

	time.Sleep(StartUpDelay * time.Second)
	slog.Debug("Delay start up to allow for the database to be ready", "delay_in_seconds", StartUpDelay)

	if err := db.MigrateDB(migrations); err != nil {
		slog.Error("Failed to apply migrations", "error", err)
		os.Exit(1)
	}

	databasePassword := os.Getenv("APP_USER_POSTGRES_PASSWORD")
	if databasePassword == "" {
		slog.Error("Missing Environment Variable: APP_USER_POSTGRES_PASSWORD")
	}

	dbInstance, err := db.NewDB(databasePassword)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbInstance.Connection.Close()

	if err := dbInstance.Connection.Ping(); err != nil {
		slog.Error("Failed to Ping the database", "error", err)
	}
	slog.Info("Connected to Database")

	if err := seedDatabase(dbInstance); err != nil {
		slog.Error("Failed to seed the database", "error", err)
	}

	users, err := dbInstance.GetAllUsers()
	if err != nil {
		slog.Error("Failed to get all users", "error", err)
	}
	slog.Debug("count of users found", "count", len(users))

	var wg sync.WaitGroup
	wg.Add(len(users))
	for _, user := range users {
		// go mainUserLoop(&wg, user)
		go mainUserLoop2(&wg, user)
	}

	wg.Wait()

	slog.Debug("Application finished running")
}
