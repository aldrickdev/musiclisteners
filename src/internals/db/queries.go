package db

import (
	"fmt"

	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/jmoiron/sqlx"
)

type QueryExecutor interface {
	GetQuery() string
	Execute(*sqlx.DB)
}

type SelectUsersResult struct {
	Users []types.User
	Err   error
}

type SelectUsers struct {
	SQL    string
	Result chan SelectUsersResult
}

func NewSelectUsers(result chan SelectUsersResult) *SelectUsers {
	return &SelectUsers{
		SQL:    SelectAllUsers,
		Result: result,
	}
}

func (q *SelectUsers) GetQuery() string {
	return q.SQL
}

func (q *SelectUsers) Execute(dbConnection *sqlx.DB) {
	allUsers := []types.User{}
	singleUser := types.User{}

	rows, err := dbConnection.Queryx(SelectAllUsers)
	if err != nil {

		q.Result <- SelectUsersResult{
			Users: allUsers,
			Err:   fmt.Errorf("Failed to query for all users: %q", err),
		}
	}
	for rows.Next() {
		err := rows.StructScan(&singleUser)
		if err != nil {
			q.Result <- SelectUsersResult{
				Users: allUsers,
				Err:   fmt.Errorf("Failed to scan for all users: %q", err),
			}
		}

		allUsers = append(allUsers, singleUser)
	}
	q.Result <- SelectUsersResult{
		Users: allUsers,
		Err:   nil,
	}
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
)
