package db

import (
	"fmt"
	"log"

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

const SelectRandomSongQuery = `
  SELECT * FROM production.available_songs
  WHERE id >= floor(random() * (SELECT max(id) FROM production.available_songs))
  ORDER BY id
  LIMIT 1;
`

const SelectSongByID = `
  SELECT * FROM production.available_songs
  WHERE id = :id
  LIMIT 1;
`

const SelectCurrentSongForUserQuery = `
  SELECT * FROM production.songs_currently_playing
  WHERE user_id = :id
  LIMIT 1;
`

const DeleteCurrentSongForUserQuery = `
  DELETE FROM production.songs_currently_playing
  WHERE user_id = :id;
`

const InsertCurrentlyPlayingSongForUserQuery = `
  INSERT INTO production.songs_currently_playing (
    user_id,
    song_id
  ) VALUES (
    :user_id,
    :song_id
  );
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

func (db *DB)SelectRandomSong() (types.Song, error) {
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

func (db *DB)SelectCurrentlyPlayingSongForUser(user types.User) (types.Song, error){
  song := types.Song{}
  currentSong := types.CurrentlyPlayingSong{}

  fmt.Printf("The user is: %v\n", user)

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

    fmt.Printf("Scanned song: %v\n", song)
  }

  return song, nil
}

func (db *DB)InsertCurrentlyPlayingSongForUser(user types.User, song types.Song) (error){
  result, err := db.Connection.NamedExec(DeleteCurrentSongForUserQuery, user)
  if err != nil {
    return fmt.Errorf("Failed to insert current song for user: %q\n", err)
  }

  rowsDeleted, err := result.RowsAffected()
  if err != nil {
    return fmt.Errorf("Failed to delete current song for user: %q\n", err)
  }
  log.Printf("%d rows were deleted\n", rowsDeleted)

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
  log.Printf("%d song inserted for the user", rowsInserted)
  return nil
}

