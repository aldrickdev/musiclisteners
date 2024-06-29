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

	// User count
	MaxUsers = 100
)

func mainUserLoop(wg *sync.WaitGroup, dbConnection *db.DB, user types.User) {
	for {
		randomSleep := rand.Intn(10) + 10
		slog.Debug("User sleeping", "user_id", user.ID, "sleep", randomSleep)
		time.Sleep(time.Second * time.Duration(randomSleep))



		randomSongChannel := make(chan db.GetRandomSongQueryResult)
		randomSongQuery := db.NewGetRandomSongQuery(randomSongChannel)
		dbConnection.QueryBuffer <- randomSongQuery

		queryBufferLength := len(dbConnection.QueryBuffer)
		slog.Debug("query buffer size", "user_id", user.ID, "queued_queries_count", queryBufferLength)
		randomSongResult := <-randomSongChannel
		slog.Debug("Got Random song", "user_id", user.ID, "song_id", randomSongResult.Song.ID)



		insertUserCurrentSongResultChannel := make(chan db.InsertUserCurrentSongResult)
		insertUserCurrentSongQuery := db.NewInsertUserCurrentSongQuery(insertUserCurrentSongResultChannel, user, randomSongResult.Song)
		dbConnection.QueryBuffer <- insertUserCurrentSongQuery

		queryBufferLength = len(dbConnection.QueryBuffer)
		slog.Debug("query buffer size", "user_id", user.ID, "queued_queries_count", queryBufferLength)
		insertUserCurrentSongResult := <-insertUserCurrentSongResultChannel
		if insertUserCurrentSongResult.Err != nil {
			slog.Error("Failed to set the current song for user", "error", insertUserCurrentSongResult.Err)
			break
		}
	}

	wg.Done()
}

func seedDatabase(dbInstance *db.DB) error {
	songs := utils.ImportCSVSongFromEmbededFS("embed/data/spotify_data.csv", data, utils.ExtractSongsFromCSVReader)
	slog.Debug("Extracted songs", "songs", songs)

	users := utils.GenerateUsers(MaxUsers)
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

	slog.Debug("Delay start up to allow for the database to be ready", "delay_in_seconds", StartUpDelay)
	time.Sleep(StartUpDelay * time.Second)

	if err := db.MigrateDB(migrations); err != nil {
		slog.Error("Failed to apply migrations", "error", err)
		os.Exit(1)
	}

	databasePassword := os.Getenv("APP_USER_POSTGRES_PASSWORD")
	if databasePassword == "" {
		slog.Error("Missing Environment Variable: APP_USER_POSTGRES_PASSWORD")
	}

	dbConnection, err := db.NewDBConnection(databasePassword, MaxUsers * 10)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbConnection.Connection.Close()

	if err := seedDatabase(dbConnection); err != nil {
		slog.Error("Failed to seed the database", "error", err)
	}

	resultChan := make(chan db.GetAllUsersQueryResult)
	query := db.NewGetAllUsersQuery(resultChan)
	dbConnection.QueryBuffer <- query

	slog.Debug("Waiting for results")
	queryResults := <-resultChan
	slog.Debug("count of users found", "count", len(queryResults.Users))

	var wg sync.WaitGroup
	wg.Add(len(queryResults.Users))
	for _, user := range queryResults.Users {
		// go mainUserLoop2(&wg, user)
		go mainUserLoop(&wg, dbConnection, user)
	}

	wg.Wait()

	slog.Debug("Application finished running")
}
