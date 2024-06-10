package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aldricdev/musiclisteners/internals/db"
	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/jmoiron/sqlx"
)

func seedDB(dbConnection *sqlx.DB, songs []types.Song, users []types.User) {
	results, err := dbConnection.NamedExec(db.InsertAvailableSong, songs)
	if err != nil {
		log.Fatalf("Failed to insert songs for seeding: %q", err)
	}

	count, err := results.RowsAffected()
	if err != nil {
		log.Printf("Driver doesn't support result type: %q\n", err)
	} else {
		fmt.Printf("Songs inserted: %v\n", count)
	}

	results, err = dbConnection.NamedExec(db.InsertUser, users)
	if err != nil {
		log.Fatalf("Failed to insert users for seeding: %q", err)
	}

	count, err = results.RowsAffected()
	if err != nil {
		log.Printf("Driver doesn't support result type: %q\n", err)
	} else {
		fmt.Printf("Users inserted: %v\n", count)
	}
}

func mainUserLoop(dbInstance db.DB, user types.User) {
	song, err := dbInstance.GetRandomSong()
	if err != nil {
		log.Printf("Failed to get a random song: %q\n", err)
	}

	fmt.Printf("Random Song: %v\n", song)
}

func main() {
	databasePassword := os.Getenv("APP_USER_POSTGRES_PASSWORD")
	if databasePassword == "" {
		log.Fatalf("Missing Environment Variable: APP_USER_POSTGRES_PASSWORD")
	}

	dbInstance := db.NewDB(databasePassword)
	defer dbInstance.Connection.Close()

	if err := dbInstance.Connection.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Database")

	// songs := utils.ImportCSVData("spotify_data.csv", utils.ExtractSongsFromCSVReader)
	// users := utils.GenerateUsers()
	// seedDB(dbInstance.Connection, songs, users)

	users, err := dbInstance.GetAllUsers()
	if err != nil {
		log.Printf("Failed to get all users: %q", err)
	}

	mainUserLoop(*dbInstance, users[0])
}
