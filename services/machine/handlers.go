package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	gin "github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context) {
	response := make(map[string]string)

	instance, err := getInstance()
	if err != nil {
		log.Printf("failed to get active instance: %v\n", err)
		response["status"] = "NOT_FOUND"
		c.JSON(http.StatusNotFound, response)
		return
	}

	response["status"] = *instance.Status
	response["redirect_link"] = fmt.Sprintf(
		"https://remotedesktop.google.com/access/session/%s",
		config.Session,
	)

	c.JSON(http.StatusOK, response)
}

// Deploys VM
func PatchHandler(c *gin.Context) {
	response := make(map[string]string)

	instance, err := getInstance()
	if err == nil {
		log.Printf("instance already running: %v\n", err)
		response["status"] = *instance.Status
		response["redirect_link"] = fmt.Sprintf(
			"https://remotedesktop.google.com/access/session/%s",
			config.Session,
		)

		c.JSON(http.StatusOK, response)
		return
	}

	addr := ipSource(c.Request)
	log.Printf("client addr is %s\n", addr)

	loc, err := ipLocation(addr)
	if err != nil {
		log.Printf("failed to locate ip: %v\n", err)
		response["status"] = "UKNOWN_SOURCE"
		c.JSON(http.StatusPreconditionFailed, response)
		return
	}

	region, distance := closestRegion(loc)
	log.Printf("closest region is %s, %.2f kilometers away\n", region, distance)

	err = deployInstance(region)
	if err != nil {
		log.Printf("failed to deploy VM: %v\n", err)
		response["status"] = "NON_DEPLOYABLE"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	instance, err = getInstance()
	if err != nil {
		log.Printf("failed to get VM status: %v\n", err)
		response["status"] = "UNKNOWN"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	for *instance.Status != "RUNNING" {
		time.Sleep(time.Second * 4)

		instance, err = getInstance()
		if err != nil {
			log.Printf("failed to get VM status: %v\n", err)
			response["status"] = "UNKNOWN"
			c.JSON(http.StatusInternalServerError, response)
			return
		}
	}

	// Instance running
	time.Sleep(time.Second * 8) // wait for boot

	response = make(map[string]string)
	response["status"] = *instance.Status
	response["redirect_link"] = fmt.Sprintf(
		"https://remotedesktop.google.com/access/session/%s",
		config.Session,
	)

	c.JSON(http.StatusOK, response)
}

func OptionsHandler(c *gin.Context) {}

// Takes disk snapshot and terminates VM
func DeleteHandler(c *gin.Context) {
	response := make(map[string]string)

	instance, err := getInstance()
	if err != nil {
		log.Printf("instance not found: %v\n", err)
		response["status"] = "MISSING"
		c.JSON(http.StatusGone, response)
		return
	}

	err = destroyInstance(instance)
	if err != nil {
		log.Printf("failed to destroy VM: %v\n", err)
		response["status"] = "NON_DESTROYABLE"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response["status"] = *instance.Status
	c.JSON(http.StatusOK, response)
}

func UserContextFromAPI(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().
		Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
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
