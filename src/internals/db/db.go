package db

import (
	"fmt"

	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const InsertAvailableSong = `
  INSERT INTO production.available_songs (
    track_name, 
    artists_name, 
    released_year
  ) VALUES (
    :track_name, 
    :artists_name, 
    :released_year
  )
`

const SelectRandomSong = `
  SELECT * FROM production.available_songs
  WHERE id >= floor(random() * (SELECT max(id) FROM production.available_songs))
  ORDER BY id
  LIMIT 1;
`

const SelectAllUsers = `
  SELECT * FROM production.users
`

const InsertUser = `
  INSERT INTO production.users (
    name,
    avatar
  ) VALUES (
    :name,
    :avatar
  )
`

type DB struct {
  Connection *sqlx.DB
}

func NewDB(password string) *DB {
	connectString := fmt.Sprintf("user=app dbname=musiclisteners sslmode=disable password=%s", password)
	db := sqlx.MustConnect("postgres", connectString)

	return &DB{
    Connection: db,
  }
}

func (db *DB)GetRandomSong() (types.Song, error) {
  randomSong := types.Song{}
  row, err := db.Connection.Queryx(SelectRandomSong)
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

func (db *DB)GetAllUsers() ([]types.User, error){
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

