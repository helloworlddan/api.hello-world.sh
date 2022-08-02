package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	gin "github.com/gin-gonic/gin"

	compute "google.golang.org/api/compute/v1"
)

//"https://remotedesktop.google.com/access/session/ca683f00-d51c-4f1a-af5e-5f9a25b3f4a8"

// Principal represents the identity that originally authorized the context of an interaction
type Principal struct {
	ID         string `header:"user_id" firestore:"user_id" json:"user_id"`
	Email      string `header:"email" firestore:"email" json:"email"`
	Name       string `header:"name" firestore:"name" json:"name"`
	PictureURL string `header:"picture" firestore:"picture" json:"picture"`
}

// Config holds the complete context of the managed machine
type Config struct {
	Session     string
	Owner       string
	Zone        string
	Project     string
	Machine     string
	Environment string
}

var (
	config *Config
)

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
	}

	router.Run()
}

func configure() (Config, error) {
	cfg := Config{}
	cfg.Session = os.Getenv("TOP_SESSION")
	cfg.Owner = os.Getenv("TOP_OWNER")
	cfg.Zone = os.Getenv("TOP_ZONE")
	cfg.Machine = os.Getenv("TOP_MACHINE")
	cfg.Project = os.Getenv("GOOGLE_CLOUD_PROJECT")
	cfg.Environment = os.Getenv("ENVIRONMENT")

	if cfg.Session == "" || cfg.Owner == "" || cfg.Zone == "" || cfg.Machine == "" || cfg.Project == "" || cfg.Environment == "" {
		return Config{}, errors.New("config incomplete")
	}
	return cfg, nil
}

func GetHandler(c *gin.Context) {
	ctx := c.Request.Context()
	gce, err := compute.NewService(ctx)
	if err != nil {
		log.Printf("failed to initialize GCE client: %v\n", err)
	}

	instances := compute.NewInstancesService(gce)
	getCall := instances.Get(config.Project, config.Zone, config.Machine)
	instance, err := getCall.Do()
	if err != nil {
		log.Printf("failed to initialize GCE client: %v\n", err)
	}

	c.JSON(http.StatusOK, instance.Status)
}

func PatchHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "unimplemented")
}

func DeleteHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "unimplemented")
}

func UserContextFromAPI(c *gin.Context) {
	var caller Principal

	if config.Environment != "dev" {
		encoded := c.Request.Header.Get("X-Endpoint-API-UserInfo")
		if encoded == "" {
			log.Printf("error: %v\n", fmt.Errorf("missing gateway user info header"))
			c.JSON(http.StatusUnauthorized, "missing gateway user info header")
			c.Abort()
			return
		}
		bytes, err := base64.RawURLEncoding.DecodeString(encoded)
		if err != nil {
			log.Printf("error: %v\n", err)
			c.JSON(http.StatusUnauthorized, "failed to decode user info header")
			c.Abort()
			return
		}

		err = json.Unmarshal(bytes, &caller)
		if err != nil {
			log.Printf("error: %v\n", err)
			c.JSON(http.StatusUnauthorized, "failed to deserialize user info header")
			c.Abort()
			return
		}

		if caller.Email != config.Owner {
			log.Println("error: unauthorized caller")
			c.JSON(http.StatusUnauthorized, "failed to deserialize user info header")
			c.Abort()
			return
		}
	} else {
		caller = Principal{
			ID:         "1",
			Name:       "dev user",
			Email:      "stay@puft.gb",
			PictureURL: "https://www.eric-andreu.com/wp-content/uploads/2021/12/account.eric-andreu-956x1024.png",
		}
	}

	// Context OK
	c.Set("principal", &caller)
	c.Next()
}
