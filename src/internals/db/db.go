package db

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type DB struct {
	Connection *sqlx.DB
}

// TODO: Add user parameter so that this can be reused in MigrateDB
func NewDB(password string) (*DB, error) {
	connectString := fmt.Sprintf("host=db user=app dbname=musiclisteners sslmode=disable password=%s", password)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return nil, err
	}

	return &DB{
		Connection: db,
	}, nil
}

func MigrateDB(migrations fs.FS) error {
	connectString := fmt.Sprintf("host=db user=postgres dbname=musiclisteners sslmode=disable password=%s", "example")
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		slog.Error("Failed to select a dialect", "error", err)
		os.Exit(1)
	}

	if err = goose.Up(db.DB, "embed/migrations"); err != nil {
		slog.Error("Failed to apply migrations", "error", err)
		os.Exit(1)
	}

	return nil
}

func (db *DB) InsertSeedStatus(status bool) error {
	statusInt := 0
	if status {
		statusInt = 1
	}

	results, err := db.Connection.NamedExec(InsertSeedStatusQuery, types.Seed{
		Status: statusInt,
	})
	if err != nil {
		slog.Error("Failed to set seed status", "status", statusInt, "error", err)
		return err
	}

	count, err := results.RowsAffected()
	if err != nil {
		slog.Warn("Failed to get count of rows affect for seed insert", "error", err)
	} else {
		slog.Debug("Row added in the seed table", "rows", count)
	}

	return nil
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
	row, err := db.Connection.Queryx(SelectRandomSongQuery)
	if err != nil {
		return types.Song{}, fmt.Errorf("Failed to get random song: %q", err)
	}

	if row.Next() {
		err = row.StructScan(&randomSong)
		if err != nil {
			return types.Song{}, fmt.Errorf("Failed to scan the random song returned: %q", err)
		}

		return randomSong, nil
	}

	return types.Song{}, fmt.Errorf("No songs returned")
}

func (db *DB) GetAllUsers() ([]types.User, error) {
	allUsers := []types.User{}
	singleUser := types.User{}
	rows, err := db.Connection.Queryx(SelectAllUsers)
	if err != nil {
		return allUsers, fmt.Errorf("Failed to query for all users: %q", err)
	}
	for rows.Next() {
		err := rows.StructScan(&singleUser)
		if err != nil {
			return allUsers, fmt.Errorf("Failed to scan for all users: %q", err)
		}

		allUsers = append(allUsers, singleUser)
	}
	return allUsers, nil
}

func (db *DB) SelectCurrentlyPlayingSongForUser(user types.User) (types.Song, error) {
	song := types.Song{}
	currentSong := types.CurrentlyPlayingSong{}

	rows, err := db.Connection.NamedQuery(SelectCurrentSongForUserQuery, user)
	if err != nil {
		return song, fmt.Errorf("Failed to query for current song: %q", err)
	}
	for rows.Next() {
		err := rows.StructScan(&currentSong)
		if err != nil {
			return song, fmt.Errorf("Failed to scan for current song: %q", err)
		}
	}

	rows, err = db.Connection.NamedQuery(SelectSongByID, map[string]any{
		"id": currentSong.SongID,
	})
	if err != nil {
		return song, fmt.Errorf("Failed to query for song: %q", err)
	}
	for rows.Next() {
		err := rows.StructScan(&song)
		if err != nil {
			return song, fmt.Errorf("Failed to scan for song: %q", err)
		}

		slog.Debug("Got current song for user", "user_id", user.ID, "song_id", song.ID)
	}

	return song, nil
}

func (db *DB) InsertCurrentlyPlayingSongForUser(user types.User, song types.Song) error {
	result, err := db.Connection.NamedExec(DeleteCurrentSongForUserQuery, user)
	if err != nil {
		return fmt.Errorf("Failed to insert current song for user: %q\n", err)
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to delete current song for user: %q\n", err)
	}
	slog.Debug("rows were deleted", "deleted_rows", rowsDeleted)

	currentSongForUser := types.CurrentlyPlayingSong{
		UserID: user.ID,
		SongID: song.ID,
	}
	result, err = db.Connection.NamedExec(InsertCurrentlyPlayingSongForUserQuery, currentSongForUser)
	if err != nil {
		return fmt.Errorf("Failed to insert current song for user: %q\n", err)
	}

	rowsInserted, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to obtain the count of rows affect: %q\n", err)
	}
	slog.Debug("Current song playing inserted for user", "count", rowsInserted, "user_id", user.ID, "song_id", song.ID)
	return nil
}

func (db *DB) InsertCurrentlyPlayingSongForUserTrans(user types.User, song types.Song) error {
	tx, err := db.Connection.Beginx()
	if err != nil {
		return fmt.Errorf("Failed to create the Transaction: %q", err)
	}
	result, err := tx.NamedExec(DeleteCurrentSongForUserQuery, user)
	if err != nil {
		return fmt.Errorf("Failed to insert current song for user: %q\n", err)
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to delete current song for user: %q\n", err)
	}
	slog.Debug("Currently playing song deleted for user", "count", rowsDeleted, "user_id", user.ID, "song_id", song.ID)

	currentSongForUser := types.CurrentlyPlayingSong{
		UserID: user.ID,
		SongID: song.ID,
	}
	result, err = tx.NamedExec(InsertCurrentlyPlayingSongForUserQuery, currentSongForUser)
	if err != nil {
		return fmt.Errorf("Failed to insert current song for user: %q\n", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to run insert current song transaction for user: %q\n", err)
	}

	rowsInserted, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to obtain the count of rows affect: %q\n", err)
	}
	slog.Debug("Currently playing song inserted for user", "count", rowsInserted, "user_id", user.ID)
	return nil
}
