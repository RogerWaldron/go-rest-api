package main

import (
	"context"
	"fmt"

	"github.com/RogerWaldron/go-rest-api/server/db"
	"github.com/RogerWaldron/go-rest-api/server/internal/comment"
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
	store, err := db.NewDatabase()
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
	
	commentService := comment.NewService(store)
	commentService.Store.PostComment(
		context.Background(),
		comment.Comment{
			ID: "123456",
			Slug: "test",
			Author: "Me",
			Body: "Works or not",
		},
	)
	fmt.Println(commentService.GetCommentByID(
		context.Background(),
		"123456",
	))

	return nil 
}

func main(){
	err := Run()
	if err != nil {
		fmt.Println(err) 
	}
}