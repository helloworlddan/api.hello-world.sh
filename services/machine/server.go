package main

import (
	"errors"
	"log"
	"os"

	gin "github.com/gin-gonic/gin"
	haversine "github.com/umahmood/haversine"
)

// Principal represents the identity that originally authorized the context of an interaction
type Principal struct {
	ID         string `header:"user_id" firestore:"user_id" json:"user_id"`
	Email      string `header:"email"   firestore:"email"   json:"email"`
	Name       string `header:"name"    firestore:"name"    json:"name"`
	PictureURL string `header:"picture" firestore:"picture" json:"picture"`
}

// Config holds the complete context of the managed machine
type Config struct {
	Session     string
	Owner       string
	Project     string
	Machine     string
	Snapshot    string
	Disk        string
	Environment string
	Regions     map[string]haversine.Coord
}

var config *Config

func configure() (Config, error) {
	cfg := Config{}
	cfg.Session = os.Getenv("TOP_SESSION")
	cfg.Owner = os.Getenv("TOP_OWNER")
	cfg.Machine = os.Getenv("TOP_MACHINE")
	cfg.Snapshot = os.Getenv("TOP_SNAPSHOT")
	cfg.Disk = os.Getenv("TOP_DISK")
	cfg.Project = os.Getenv("GOOGLE_CLOUD_PROJECT")
	cfg.Environment = os.Getenv("ENVIRONMENT")
	cfg.Regions = seedRegions()

	if cfg.Session == "" || cfg.Owner == "" || cfg.Machine == "" ||
		cfg.Snapshot == "" ||
		cfg.Disk == "" ||
		cfg.Project == "" ||
		cfg.Environment == "" {
		return Config{}, errors.New("config incomplete")
	}
	return cfg, nil
}

func main() {
	cfg, err := configure()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	config = &cfg

	router := gin.Default()
	machine := router.Group("/machine")
	machine.Use(UserContextFromAPI)
	{
		machine.GET("/", GetHandler)
		machine.PATCH("/", PatchHandler)
		machine.DELETE("/", DeleteHandler)
		machine.OPTIONS("/", OptionsHandler)
	}

	router.Run()
}
