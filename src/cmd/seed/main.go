package main

import (
	"log"
	"os"

	"github.com/aldricdev/musiclisteners/internals/db"
	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/aldricdev/musiclisteners/internals/utils"
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
		log.Printf("Songs inserted: %v\n", count)
	}

	results, err = dbConnection.NamedExec(db.InsertUser, users)
	if err != nil {
		log.Fatalf("Failed to insert users for seeding: %q", err)
	}

	count, err = results.RowsAffected()
	if err != nil {
		log.Printf("Driver doesn't support result type: %q\n", err)
	} else {
		log.Printf("Users inserted: %v\n", count)
	}
}

func main() {
	databasePassword := os.Getenv("APP_USER_POSTGRES_PASSWORD")
	if databasePassword == "" {
		log.Fatalf("Missing Environment Variable: APP_USER_POSTGRES_PASSWORD")
	}

	dbInstance, err := db.NewDB(databasePassword)
	if err != nil {
		log.Fatalf("Failed to connect to database: %q", err)
	}
	defer dbInstance.Connection.Close()

	if err := dbInstance.Connection.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Database")

	songs := utils.ImportCSVData("spotify_data.csv", utils.ExtractSongsFromCSVReader)
	users := utils.GenerateUsers()
	seedDB(dbInstance.Connection, songs, users)
}
