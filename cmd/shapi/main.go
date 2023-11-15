package main

import (
	"encoding/json"
	"log"
	"os"

	"apachejuice.dev/apachejuice/shopping-list-api/internal/api"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/logging"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/repo"
)

type Config struct {
	LogFile string                  `json:"logfile"`
	Auth    api.AuthenticatorConfig `json:"auth"`
}

func main() {
	config, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	var cnf Config
	err = json.Unmarshal(config, &cnf)
	if err != nil {
		log.Fatal(err)
	}

	logging.InitLog(cnf.LogFile)
	repo.Connect()

	api := api.NewApiImpl(api.NewAuthenticator(cnf.Auth))
	api.Run("localhost:9099")
}
