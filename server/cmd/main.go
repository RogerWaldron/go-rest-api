package main

import (
	"context"
	"fmt"
	"log"

	"github.com/RogerWaldron/go-rest-api/server/db"
	"github.com/joho/godotenv"
)

func Run() error {
	fmt.Println("starting app...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db, err := db.NewDatabase()
	if err != nil { 
		return err
	} 

	err = db.Ping(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("database connected oks")
	
	return nil 
}
func main(){
	err := Run()
	if err != nil {
		fmt.Println(err) 
	}
}