package db

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/aldricdev/musiclisteners/internals/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type DB struct {
	Connection  *sqlx.DB
	QueryBuffer chan QueryExecutor
}

func NewDBConnection(user string, password string, queryBufferSize int) (*DB, error) {
	connectString := fmt.Sprintf("host=db user=%s dbname=musiclisteners sslmode=disable password=%s", user, password)
	connection, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return nil, err
	}

	db := &DB{
		Connection:  connection,
		QueryBuffer: make(chan QueryExecutor, queryBufferSize),
	}

	go db.connectionLoop()

	return db, nil
}

func MigrateDB(migrations fs.FS) error {
	postgresUserPassword := utils.MustGetEnv("POSTGRES_PASSWORD")

	db, err := NewDBConnection("postgres", postgresUserPassword, 5)
	if err != nil {
		return err
	}
	defer db.Connection.Close()

	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		slog.Error("Failed to select a dialect", "error", err)
		os.Exit(1)
	}

	if err = goose.Up(db.Connection.DB, "embed/migrations"); err != nil {
		slog.Error("Failed to apply migrations", "error", err)
		os.Exit(1)
	}

	return nil
}

func (db *DB) SeedDB(songs []types.Song, users []types.User) error {
	err := db.InsertAvailableSongBatch(songs)
	if err != nil {
		slog.Error("Failed to insert a batch of available songs", "error", err)
		return err
	}

	err = db.InsertUserBatch(users)
	if err != nil {
		slog.Error("Failed to insert a batch of users", "error", err)
		return err
	}

	return nil
}

func (db *DB) connectionLoop() {
	for q := range db.QueryBuffer {
		slog.Debug("Query received", "query", q.GetQuery())
		db.connectionHandler(q)
	}
}

func (db *DB) connectionHandler(query QueryExecutor) {
	query.Execute(db.Connection)
}

func (db *DB) InsertAvailableSongBatch(songs []types.Song) error {
	results, err := db.Connection.NamedExec(InsertAvailableSongQuery, songs)
	if err != nil {
		slog.Error("Failed to seed database with songs", "error", err)
		return err
	}

	count, err := results.RowsAffected()
	if err != nil {
		slog.Warn("Driver doesn't support result type", "error", err)
	} else {
		slog.Debug("Amount of songs inserted", "count", count)
	}

	return nil
}

func (db *DB) InsertUserBatch(users []types.User) error {
	results, err := db.Connection.NamedExec(InsertUserQuery, users)
	if err != nil {
		slog.Error("Failed to seed database with users", "error", err)
		return err
	}

	count, err := results.RowsAffected()
	if err != nil {
		slog.Warn("Driver doesn't support result type", "error", err)
	} else {
		slog.Debug("Amount of users inserted", "count", count)
	}

	return nil
}

func (db *DB) SelectRandomSong() (types.Song, error) {
	randomSong := types.Song{}

	randomSongChannel := make(chan GetRandomSongQueryResult)
	randomSongQuery := NewGetRandomSongQuery(randomSongChannel)
	db.QueryBuffer <- randomSongQuery

	queryBufferLength := len(db.QueryBuffer)
	slog.Debug("query buffer size", "queued_queries_count", queryBufferLength)

	randomSongResult := <-randomSongChannel
	if randomSongResult.Err != nil {
		return randomSong, randomSongResult.Err
	}

	return randomSongResult.Song, nil
}

func (db *DB) SelectAllUsers() ([]types.User, error) {
	resultChan := make(chan GetAllUsersQueryResult)
	query := NewGetAllUsersQuery(resultChan)
	db.QueryBuffer <- query

	slog.Debug("Waiting for results")
	queryResults := <-resultChan
	if queryResults.Err != nil {
		return queryResults.Users, queryResults.Err
	}
	slog.Debug("count of users found", "count", len(queryResults.Users))

	return queryResults.Users, nil
}

func (db *DB) InsertCurrentlyPlayingSongForUser(user types.User, song types.Song) error {
	insertUserCurrentSongResultChannel := make(chan InsertUserCurrentSongResult)
	insertUserCurrentSongQuery := NewInsertUserCurrentSongQuery(insertUserCurrentSongResultChannel, user, song)
	db.QueryBuffer <- insertUserCurrentSongQuery

	queryBufferLength := len(db.QueryBuffer)
	slog.Debug("query buffer size", "user_id", user.ID, "queued_queries_count", queryBufferLength)

	insertUserCurrentSongResult := <-insertUserCurrentSongResultChannel
	if insertUserCurrentSongResult.Err != nil {
		return insertUserCurrentSongResult.Err
	}

	return nil
}
