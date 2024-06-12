package types

type Song struct {
	ID           int    `db:"id"`
	Name         string `db:"track_name"`
	Artists      string `db:"artists_name"`
	ReleasedYear int    `db:"released_year"`
}

type User struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Avatar string `db:"avatar"`
}

type CurrentlyPlayingSong struct {
	ID     int `db:"id"`
	UserID int `db:"user_id"`
	SongID int `db:"song_id"`
}
