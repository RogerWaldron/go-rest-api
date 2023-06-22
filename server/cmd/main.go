package main

import (
	"context"
	"fmt"
	"os"

	"github.com/RogerWaldron/go-rest-api/server/db"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func Run() error {
	fmt.Println("starting app...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg(".env failed to load file")
		return err
	}
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_TABLE"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)
	store, err := db.NewDatabase(connectionString)
	if err != nil { 
		log.Fatal().Err(err).Msg("")
		return err
	} 

	err = store.Ping(context.Background())
	if err != nil {
		return err
	}

	err = store.MigrateDB()
	if err != nil {
		log.Fatal().Err(err).Msg("migration failed to setup database")
		return err
	}

	return nil 
}

func main(){
	err := Run()
	if err != nil {
		fmt.Println(err) 
	}
}