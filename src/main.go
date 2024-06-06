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
)

type Song struct {
	Name string
	Artists string
	ReleasedYear int
}

// CSv Shape
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

func main() {
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

		jsonB, err := json.Marshal(Song{
				Name: row[TrackName],
				Artists: row[ArtistsNames],
				ReleasedYear: releaseYear,
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stdout, "%s\n", jsonB)
	}

}
