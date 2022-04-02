package main

import (
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/sergeysynergy/gopracticum/agent"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	err := godotenv.Load("./config/.env")
	if err != nil {
		err = godotenv.Load("../../config/.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {
	port := os.Getenv("SERVER_PORT")

	cfg := agent.Config{
		PollInterval:   1 * time.Second, // in prod 2
		ReportInterval: 2 * time.Second, // in prod 10
		Port:           port,
	}

	agent, err := agent.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	agent.Run()
}
