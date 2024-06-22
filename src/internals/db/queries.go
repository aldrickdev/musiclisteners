package db

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
