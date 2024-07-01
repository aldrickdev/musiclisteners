package db

import (
	"fmt"
	"log/slog"

	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/jmoiron/sqlx"
)

type QueryExecutor interface {
	GetQuery() []string
	Execute(*sqlx.DB)
}

type InsertUserCurrentSongResult struct {
	Err error
}

type InsertUserCurrentSongQuery struct {
	SQL    []string
	User   types.User
	Song   types.Song
	Result chan InsertUserCurrentSongResult
}

func NewInsertUserCurrentSongQuery(result chan InsertUserCurrentSongResult, user types.User, song types.Song) *InsertUserCurrentSongQuery {
	return &InsertUserCurrentSongQuery{
		SQL:    []string{DeleteCurrentSongForUserQuery, InsertCurrentlyPlayingSongForUserQuery, InsertSongHistory},
		User:   user,
		Song:   song,
		Result: result,
	}
}

func (q *InsertUserCurrentSongQuery) GetQuery() []string {
	return q.SQL
}

func (q *InsertUserCurrentSongQuery) Execute(dbConnection *sqlx.DB) {
	tx, err := dbConnection.Beginx()
	if err != nil {
		q.Result <- InsertUserCurrentSongResult{
			Err: fmt.Errorf("Failed to create the Transaction: %q", err),
		}
		return
	}

	// Delete users currently playing song
	result, err := tx.NamedExec(q.SQL[0], q.User)
	if err != nil {

		q.Result <- InsertUserCurrentSongResult{
			Err: fmt.Errorf("Failed to insert current song for user: %q", err),
		}
		return
	}
	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		slog.Debug("Checking count of rows affected is not supported", "error", err)
	}
	slog.Debug("Currently playing song deleted for user", "rows_deleted", rowsDeleted, "user_id", q.User.ID)


	// Insert users currently playing song
	currentSongForUser := types.CurrentlyPlayingSong{
		UserID: q.User.ID,
		SongID: q.Song.ID,
	}
	result, err = tx.NamedExec(q.SQL[1], currentSongForUser)
	if err != nil {
		q.Result <- InsertUserCurrentSongResult{
			Err: fmt.Errorf("Failed to insert current song for user: %q", err),
		}
		return
	}
	rowsInserted, err := result.RowsAffected()
	if err != nil {
		slog.Debug("Checking count of rows affected is not supported", "error", err)
	}
	slog.Debug("Currently playing song inserted for user", "rows_inserted", rowsInserted, "user_id", q.User.ID)


	// Insert users song into song history
	songHistory := types.SongHistory{
		UserID: q.User.ID,
		SongID: q.Song.ID,
	}
	result, err = tx.NamedExec(q.SQL[2], songHistory)
	if err != nil {
		q.Result <- InsertUserCurrentSongResult{
			Err: fmt.Errorf("Failed to insert song history for user: %q", err),
		}
		return
	}
	rowsInserted, err = result.RowsAffected()
	if err != nil {
		slog.Debug("Checking count of rows affected is not supported", "error", err)
	}
	slog.Debug("Song inserted into song history", "rows_inserted", rowsInserted, "user_id", q.User.ID, "song_id", q.Song.ID)


	err = tx.Commit()
	if err != nil {
		q.Result <- InsertUserCurrentSongResult{
			Err: fmt.Errorf("Failed to commit insert current song transaction for user: %q", err),
		}
		return
	}

	q.Result <- InsertUserCurrentSongResult{
		Err: nil,
	}

	slog.Info("Currently playing song inserted for user", "inserted_count", rowsInserted, "user_id", q.User.ID, "song_id", q.Song.ID)
}

const (
	InsertAvailableSongQuery = `
		INSERT INTO production.available_songs (
			track_name, 
			artists_name, 
			released_year
		) VALUES (
			:track_name, 
			:artists_name, 
			:released_year
		);
	`

	InsertSeedStatusQuery = `
		INSERT INTO production.seed (
			status
		) VALUES (
			:status
		);
	`

	SelectSeedStatusQuery = `
		SELECT * FROM production.seed;
	`

	SelectRandomSongQuery = `
		SELECT * FROM production.available_songs
		WHERE id >= floor(random() * (SELECT max(id) FROM production.available_songs))
		ORDER BY id
		LIMIT 1;
	`

	SelectSongByID = `
		SELECT * FROM production.available_songs
		WHERE id = :id
		LIMIT 1;
	`

	SelectCurrentSongForUserQuery = `
		SELECT * FROM production.songs_currently_playing
		WHERE user_id = :id
		LIMIT 1;
	`

	DeleteCurrentSongForUserQuery = `
		DELETE FROM production.songs_currently_playing
		WHERE user_id = :id;
	`

	InsertCurrentlyPlayingSongForUserQuery = `
		INSERT INTO production.songs_currently_playing (
			user_id,
			song_id
		) VALUES (
			:user_id,
			:song_id
		);
	`

	SelectAllUsers = `
		SELECT * FROM production.users;
	`

	InsertUserQuery = `
		INSERT INTO production.users (
			name,
			avatar
		) VALUES (
			:name,
			:avatar
		);
	`

	InsertSongHistory = `
		INSERT INTO production.song_history (
			user_id,
			song_id
		) VALUES (
			:user_id,
			:song_id
		);
	`
)
