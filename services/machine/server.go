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
	Session string
	Owner   string
	Zone    string
	Project string
	Machine string
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

	if cfg.Session == "" || cfg.Owner == "" || cfg.Zone == "" || cfg.Machine == "" || cfg.Project == "" {
		return Config{}, errors.New("config incomplete")
	}
	return cfg, nil
}

func GetHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "unimplemented")
}

func PatchHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "unimplemented")
}

func DeleteHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "unimplemented")
}

func UserContextFromAPI(c *gin.Context) {
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

	var caller Principal
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

	// Context OK
	c.Set("principal", &caller)
	c.Next()
}
