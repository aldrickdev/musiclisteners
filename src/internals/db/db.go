package db

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
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

