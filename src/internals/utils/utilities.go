package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/go-faker/faker/v4"
)

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

// Using generics incase other data will be provided using CSV
func ImportCSVData[T types.Song](csvFile string, handleRead func(*csv.Reader) []T) []T {
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open file: %q", err)
	}
	csvReader := csv.NewReader(bufio.NewReader(file))
	return handleRead(csvReader)
}

type FakerUser struct {
  Name string `faker:"first_name"`
}

func GenerateUsers() []types.User{
  users := []types.User{}

  for i := range(20) {
    user := FakerUser{}

    err := faker.FakeData(&user)
    if err != nil{
    log.Fatalf("Failed to generate user: %q", err)
    }

    users = append(users, types.User{
    	Name:   user.Name,
      Avatar: fmt.Sprintf("http://i.pravatar.cc/149?u=%d",i),
    })
  }

  return users
}

func ExtractSongsFromCSVReader(reader *csv.Reader) []types.Song {
	availableSongs := []types.Song{}

	// get rid of header
	_, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	rowNo := 0
	for {
		rowNo++
		row, err := reader.Read()
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

		availableSongs = append(availableSongs, types.Song{
				Name:         row[TrackName],
				Artists:      row[ArtistsNames],
				ReleasedYear: releaseYear,
			},
		)
	}

	return availableSongs
}
