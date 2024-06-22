package utils

import (
	"bufio"
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/aldricdev/musiclisteners/internals/db"
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
		slog.Error("Failed to open file", "error", err)
		os.Exit(1)
	}
	csvReader := csv.NewReader(bufio.NewReader(file))
	return handleRead(csvReader)
}

func ImportCSVSongFromEmbededFS(csvFilename string, fileSystem embed.FS, handleRead func(*csv.Reader) []types.Song) []types.Song {
	csvFile, err := fileSystem.Open(csvFilename)
	if err != nil {
		slog.Error("Failed to open CSV", "filename", csvFilename, "error", err.Error())
	}

	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	return handleRead(csvReader)
}

type FakerUser struct {
	Name string `faker:"first_name"`
}

func GenerateUsers(count int) []types.User {
	users := []types.User{}

	for i := range count {
		user := FakerUser{}

		err := faker.FakeData(&user)
		if err != nil {
			log.Fatalf("Failed to generate user: %q", err)
		}

		users = append(users, types.User{
			Name:   user.Name,
			Avatar: fmt.Sprintf("http://i.pravatar.cc/149?u=%d", i),
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
			slog.Error("Failed to read next line", "error", err)
			os.Exit(1)
		}

		releaseYear, err := strconv.Atoi(row[ReleasedYear])
		if err != nil {
			slog.Warn("Failed to convert string year to integer, skipping song", "row_number", rowNo, "released_year", row[ReleasedYear], "song_name", row[TrackName])
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

func SeedDB(dbConnection *db.DB, songs []types.Song, users []types.User) error {
	err := dbConnection.InsertAvailableSongBatch(songs)
	if err != nil {
		slog.Error("Failed to insert a batch of available songs", "error", err)
		return err
	}

	err = dbConnection.InsertUserBatch(users)
	if err != nil {
		slog.Error("Failed to insert a batch of users", "error", err)
		return err
	}

	return nil
}
