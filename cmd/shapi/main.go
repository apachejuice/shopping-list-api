package main

import (
	"encoding/json"
	"log"
	"os"

	"apachejuice.dev/apachejuice/shopping-list-api/internal/api"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/logging"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/repo"
	"github.com/gin-gonic/gin"
)

type Config struct {
	LogFile string                  `json:"logfile"`
	Auth    api.AuthenticatorConfig `json:"auth"`
	Api     struct {
		Host           string   `json:"host"`
		TrustedProxies []string `json:"trusted_proxies"`
	}
}

func main() {
	if os.Getenv("DEV") != "" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

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
	api.Run(cnf.Api.Host, cnf.Api.TrustedProxies)
}
