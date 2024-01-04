package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/yves1da2vinci3/go_elastic/db"
	"github.com/yves1da2vinci3/go_elastic/handler"
)

func main() {
	var dbPort int
	var err error
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	// Load env variables
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("POSTGRES_PORT")
	fmt.Println("port :", port)
	if dbPort, err = strconv.Atoi(port); err != nil {
		logger.Err(err).Msg("failed to parse database port")
		os.Exit(1)
	}
	dbConfig := db.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     dbPort,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DB"),
		Logger:   logger,
	}
	logger.Info().Interface("config", &dbConfig).Msg("config:")
	dbInstance, err := db.Init(dbConfig)
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}
	logger.Info().Msg("Database connection established")

	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}

	h := handler.New(dbInstance, esClient, logger)
	router := gin.Default()
	rg := router.Group("/v1")
	h.Register(rg)
	router.Run(":8080")
}
