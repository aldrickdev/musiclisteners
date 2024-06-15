package main

import (
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/aldricdev/musiclisteners/internals/db"
	"github.com/aldricdev/musiclisteners/internals/types"
)

func mainUserLoop(wg *sync.WaitGroup, user types.User) {
	databasePassword := os.Getenv("APP_USER_POSTGRES_PASSWORD")
	if databasePassword == "" {
		log.Fatalf("Missing Environment Variable: APP_USER_POSTGRES_PASSWORD")
	}

	for {
		randomSleep := rand.Intn(10) + 0
		log.Printf("User %d sleeping for %d seconds\n", user.ID, randomSleep)

		dbInstance, err := db.NewDB(databasePassword)
		if err != nil {
			log.Printf("Failed to connect to database: %q", err)
			continue
		}

		time.Sleep(time.Duration(randomSleep) * time.Second)

		song, err := dbInstance.SelectRandomSong()
		if err != nil {
			log.Print(err.Error())
			break
		}

		if err = dbInstance.InsertCurrentlyPlayingSongForUserTrans(user, song); err != nil {
			log.Print(err)
			continue
		}

		song, err = dbInstance.SelectCurrentlyPlayingSongForUser(user)
		if err != nil {
			log.Printf("Failed to get the current song for user: %v, error: %q\n", user.ID, err)
			break
		}

		log.Printf("Song currently playing: %s\n", song.Name)
		dbInstance.Connection.Close()
	}

	wg.Done()
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
		log.Fatalf("Failed to Ping the database: %q", err)
	}
	log.Println("Connected to Database")

	users, err := dbInstance.GetAllUsers()
	if err != nil {
		log.Fatalf("Failed to get all users: %q", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(users))
	for _, user := range users {
		go mainUserLoop(&wg, user)
	}
	wg.Wait()
}
