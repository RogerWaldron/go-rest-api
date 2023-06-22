package main

import (
	"context"
	"fmt"
	"os"

	"github.com/RogerWaldron/go-rest-api/server/db"
	"github.com/RogerWaldron/go-rest-api/server/internal/comment"
	"github.com/RogerWaldron/go-rest-api/server/internal/webHTTP"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func buildConnectionString() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg(".env failed to load file")
		return "", err
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

	return connectionString, nil
}

func Run() error {
	fmt.Println("starting app...")
  connect, err := buildConnectionString()
	if err != nil {
		panic("failed to build database connection string")
	}

	db, err := db.NewDatabase(connect)
	if err != nil { 
		log.Fatal().Err(err).Msg("")
		return err
	} 

	err = db.Ping(context.Background())
	if err != nil {
		return err
	}

	err = db.MigrateDB()
	if err != nil {
		log.Fatal().Err(err).Msg("migration failed to setup database")
		return err
	}

	cmtService := comment.NewService(db)

	httpHandler := webHTTP.NewHandler(cmtService)
	err = httpHandler.Serve()
	if err != nil {
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