package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
		machine.OPTIONS("/", OptionsHandler)
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
		log.Printf("failed to retrieve instance status: %v\n", err)
	}

	response := make(map[string]string)
	response["message"] = "machine found"
	response["status"] = instance.Status

	c.JSON(http.StatusOK, response)
}

func PatchHandler(c *gin.Context) {
	ctx := c.Request.Context()
	gce, err := compute.NewService(ctx)
	if err != nil {
		log.Printf("failed to initialize GCE client: %v\n", err)
	}

	instances := compute.NewInstancesService(gce)
	startCall := instances.Start(config.Project, config.Zone, config.Machine)
	_, err = startCall.Do()
	if err != nil {
		log.Printf("failed to start instance: %v\n", err)
	}

	getCall := instances.Get(config.Project, config.Zone, config.Machine)
	instance, err := getCall.Do()
	if err != nil {
		log.Printf("failed to retrieve instance status: %v\n", err)
	}

	for instance.Status != "RUNNING" {
		time.Sleep(time.Second * 4)

		instance, err = getCall.Do()
		if err != nil {
			log.Printf("failed to retrieve instance status: %v\n", err)
		}
	}
	// Instance running
	time.Sleep(time.Second * 8) // wait for boot

	response := make(map[string]string)
	response["message"] = "boot complete"
	response["redirect_link"] = fmt.Sprintf("https://remotedesktop.google.com/access/session/%s", config.Session)

	c.JSON(http.StatusOK, response)
}

func OptionsHandler(c *gin.Context) {}

func DeleteHandler(c *gin.Context) {
	ctx := c.Request.Context()
	gce, err := compute.NewService(ctx)
	if err != nil {
		log.Printf("failed to initialize GCE client: %v\n", err)
	}

	instances := compute.NewInstancesService(gce)
	stopCall := instances.Stop(config.Project, config.Zone, config.Machine)
	_, err = stopCall.Do()
	if err != nil {
		log.Printf("failed to stop instance: %v\n", err)
	}

	response := make(map[string]string)
	response["message"] = "shutdown in progress"

	c.JSON(http.StatusOK, response)
}

func UserContextFromAPI(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "PATCH, OPTIONS, GET, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

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
