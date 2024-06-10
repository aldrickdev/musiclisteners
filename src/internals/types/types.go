package types

type Song struct {
	Name         string `db:"track_name"`
	Artists      string `db:"artists_name"`
	ReleasedYear int    `db:"released_year"`
}

type User struct {
	Name   string `db:"name"`
	Avatar string `db:"avatar"`
}
