package main

import (
	"fmt"
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
		log.Fatalf("Failed to insert songs for seeding: %q",err)
	}

	count, err := results.RowsAffected()
	if err != nil {
		log.Printf("Driver doesn't support result type: %q\n", err)
	} else {
		fmt.Printf("Songs inserted: %v\n", count)
	}

	results, err = dbConnection.NamedExec(db.InsertUser, users)
	if err != nil {
		log.Fatalf("Failed to insert users for seeding: %q",err)
	}

	count, err = results.RowsAffected()
	if err != nil {
		log.Printf("Driver doesn't support result type: %q\n", err)
	} else {
		fmt.Printf("Users inserted: %v\n", count)
	}
}

func main() {
	databasePassword := os.Getenv("APP_USER_POSTGRES_PASSWORD")
	if databasePassword == "" {
		log.Fatalf("Missing Environment Variable: APP_USER_POSTGRES_PASSWORD")
	}

	db := db.NewDB(databasePassword)
	defer db.Connection.Close()

	if err := db.Connection.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Database")

	songs := utils.ImportCSVData("spotify_data.csv", utils.ExtractSongsFromCSVReader)
	// fmt.Println(songs)
	users := utils.GenerateUsers()
	// fmt.Println(users)

	seedDB(db.Connection, songs, users)
}
