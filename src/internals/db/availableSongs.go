package db

import (
	"fmt"

	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/jmoiron/sqlx"
)

type GetRandomSongQueryResult struct {
	Song types.Song
	Err   error
}

type GetRandomSongQuery struct {
	SQL    []string
	Result chan GetRandomSongQueryResult
}

func NewGetRandomSongQuery(result chan GetRandomSongQueryResult) *GetRandomSongQuery {
	return &GetRandomSongQuery{
		SQL:    []string{SelectRandomSongQuery},
		Result: result,
	}
}

func (q *GetRandomSongQuery) GetQuery() []string {
	return q.SQL
}

func (q *GetRandomSongQuery) Execute(dbConnection *sqlx.DB) {
	randomSong := types.Song{}

	row, err := dbConnection.Queryx(q.SQL[0])
	if err != nil {
		q.Result <- GetRandomSongQueryResult{
			Song: randomSong,
			Err:   fmt.Errorf("Failed to get random song: %q", err),
		}
		return
	}

	defer row.Close()

	for row.Next() {
		err = row.StructScan(&randomSong)
		if err != nil {
			q.Result <- GetRandomSongQueryResult{
				Song: randomSong,
				Err:   fmt.Errorf("Failed to scan the random song returned: %q", err),
			}
			return
		}
	}

	q.Result <- GetRandomSongQueryResult{
		Song: randomSong,
		Err:   nil,
	}
}
