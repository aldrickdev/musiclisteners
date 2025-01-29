package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
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
	StartUpDelay = 5

	// User count
	MaxUsers = 100
)

func init() {
	// Configure the default logger
	logOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	logLevel := os.Getenv("APP_LOG_LEVEL")
	switch ll := strings.ToLower(logLevel); ll {
	case "debug":
		logOpts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}

	case "error":
		logOpts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelError,
		}
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, logOpts))
	slog.SetDefault(log)
}

func main() {
	err := tracer.Start(tracer.WithAgentAddr("datadog-agent:8126"))

	if err != nil {
		slog.Error("Failed to start tracer", "error", err)
		os.Exit(1)
	}
	defer tracer.Stop()

	slog.Debug("Delay start up to allow for the database to be ready", "delay_in_seconds", StartUpDelay)
	time.Sleep(StartUpDelay * time.Second)

	// Perform Database Migrations
	if err := db.MigrateDB(migrations); err != nil {
		slog.Error("Failed to apply migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("Migrations successfully applied")

	// Get required environment variables
	databaseUser := utils.MustGetEnv("APP_USER_POSTGRES_USERNAME")
	databasePassword := utils.MustGetEnv("APP_USER_POSTGRES_PASSWORD")

	// Get database connection
	dbConnection, err := db.NewDBConnection(databaseUser, databasePassword, MaxUsers*10)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbConnection.Connection.Close()

	// Seed the database
	if err := seedDatabase(dbConnection); err != nil {
		slog.Error("Failed to seed the database", "error", err)
	}
	slog.Info("Databases seeded")

	// Get all of the users from the database
	users, err := dbConnection.SelectAllUsers()
	if err != nil {
		slog.Error("Failed to get all users", "error", err)
		os.Exit(1)
	}
	userCount := len(users)
	slog.Debug("count of users found", "count", userCount)

	// Start the users go routines
	var wg sync.WaitGroup
	wg.Add(userCount)
	for _, user := range users {
		go mainUserLoop(&wg, dbConnection, user)
	}

	wg.Wait()
	slog.Info("Application finished running")
}

func mainUserLoop(wg *sync.WaitGroup, dbConnection *db.DB, user types.User) {

	slog.Debug("Starting main user loop", "user_id", user.ID)

	for {
		// span := tracer.StartSpan("user.mainloop", tracer.ResourceName("mainloop"))
		span, ctx := tracer.StartSpanFromContext(context.Background(), "user.mainloop", tracer.ResourceName("mainloop"))
		span.SetUser(fmt.Sprintf("%d", user.ID))

		// Get a random song
		song, err := dbConnection.SelectRandomSong(ctx)
		if err != nil {
			slog.Error("Failed to get random song for user", "user_id", user.ID, "error", err)
			break
		}

		// Set the random song as the users current song
		if err = dbConnection.InsertCurrentlyPlayingSongForUser(user, song); err != nil {
			slog.Error("Failed to insert the currently playing song for user", "user_id", user.ID, "error", err)
			break
		}
		slog.Info("Set the users current song", "user_id", user.ID, "song_id", song.ID)

		// Random sleep
		randomSleep := rand.Intn(10) + 10
		slog.Debug("User sleeping", "user_id", user.ID, "sleep", randomSleep)
		time.Sleep(time.Second * time.Duration(randomSleep))

		span.Finish()
	}

	wg.Done()
}

func seedDatabase(dbInstance *db.DB) error {
	songs := utils.ImportCSVSongFromEmbededFS("embed/data/spotify_data.csv", data, utils.ExtractSongsFromCSVReader)
	slog.Debug("Extracted songs", "songs", songs)

	users := utils.GenerateUsers(MaxUsers)
	slog.Debug("generated users", "users", users)

	if err := dbInstance.SeedDB(songs, users); err != nil {
		slog.Error("Failed to seed database", "error", err)
		return err
	}

	return nil
}
