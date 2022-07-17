package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	gin "github.com/gin-gonic/gin"
)

// Principal represents the the identity that originally authorized the context of an interaction
type Principal struct {
	ID         string `header:"user_id" firestore:"user_id" json:"user_id"`
	Email      string `header:"email" firestore:"email" json:"email"`
	Name       string `header:"name" firestore:"name" json:"name"`
	PictureURL string `header:"picture" firestore:"picture" json:"picture"`
}

func main() {
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

	// Context OK
	c.Set("principal", &caller)
	c.Next()
}
