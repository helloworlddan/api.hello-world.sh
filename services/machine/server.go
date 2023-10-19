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
	Environment string
	Regions     map[string]haversine.Coord

	InstanceName      string
	MachineType       string
	Preemtibility     bool
	ProvisioningModel string
	AutoRestart       bool
	OnHostMaintenance string
	SnapshotName      string
	DiskName          string
	DiskType          string
}

var config *Config

func configure() (Config, error) {
	cfg := Config{}
	cfg.Session = os.Getenv("TOP_SESSION")
	cfg.Owner = os.Getenv("TOP_OWNER")
	cfg.Project = os.Getenv("GOOGLE_CLOUD_PROJECT")
	cfg.Environment = os.Getenv("ENVIRONMENT")
	cfg.Regions = seedRegions()

	if cfg.Session == "" ||
		cfg.Owner == "" ||
		cfg.Project == "" ||
		cfg.Environment == "" {
		return Config{}, errors.New("config incomplete")
	}

	// Runs defaults
	cfg.InstanceName = "top"
	cfg.MachineType = "n2d-highcpu-4"
	cfg.Preemtibility = true
	cfg.ProvisioningModel = "SPOT"
	cfg.AutoRestart = true
	cfg.OnHostMaintenance = "TERMINATE"
	cfg.SnapshotName = "top"
	cfg.DiskName = "top"
	cfg.DiskType = "pd-ssd"

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
