package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Song struct {
	Name         string
	Artists      string
	ReleasedYear int
}

// CSV Shape
const (
	TrackName = iota
	ArtistsNames
	ArtistCount
	ReleasedYear
	ReleasedMonth
	ReleasedDay
	SpotifyPlaylistCount
	SpotifyPlaylistChart
	NumberOfStreams
	ApplePlaylistCount
	AppleChartCount
	DeezerPlaylistCount
	DeezerChartCount
	ShazamChartCount
	BPM
	Key
	Mode
	Danceability
	Valence
	Energy
	Acousticness
	Instrumentalness
	Liveness
	Speechiness
)

func getDB(password string) *sqlx.DB {
	connectString := fmt.Sprintf("user=app dbname=musiclisteners sslmode=disable password=%s", password)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func main() {
	databasePassword := os.Getenv("APP_USER_POSTGRES_PASSWORD")
	if databasePassword == "" {
		log.Fatalf("Missing Environment Variable: APP_USER_POSTGRES_PASSWORD")
	}

	db := getDB("yoo")

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Database")

	file, err := os.Open("spotify_data.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %q", err)
	}

	csvReader := csv.NewReader(bufio.NewReader(file))

	// get rid of header
	_, err = csvReader.Read()
	if err != nil {
		log.Fatal(err)
	}

	rowNo := 0
	for {
		rowNo++
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		releaseYear, err := strconv.Atoi(row[ReleasedYear])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%d:%d strconv.Atoi(%s): %+v\n", rowNo, ReleasedYear, row[ReleasedYear], err)
			continue
		}

		song := Song{
			Name:         row[TrackName],
			Artists:      row[ArtistsNames],
			ReleasedYear: releaseYear,
		}

		_, err = db.NamedExec(
			`INSERT INTO production.available_songs (track_name, artists_name, released_year) VALUES (:track_name, :artists_name, :released_year)`,
			map[string]interface{}{
				"track_name":    song.Name,
				"artists_name":  song.Artists,
				"released_year": song.ReleasedYear,
			},
		)
		if err != nil {
			fmt.Printf("Failed to insert the data: %v+, error: %q\n", song, err)
		}

		_, err = json.Marshal(song)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Fprintf(os.Stdout, "%s\n", jsonB)
	}
}
