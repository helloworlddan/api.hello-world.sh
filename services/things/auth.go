package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	auth "firebase.google.com/go/v4/auth"
	gin "github.com/gin-gonic/gin"
)

// Caller represents the calling API user
type Caller struct {
	ID         string `json:"user_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	PictureURL string `json:"picture"`
}

// Authenticate implements the security middleware
func Authenticate(c *gin.Context) {
	ctx := c.Request.Context()
	c.Set("trace.context", ctx)

	if Global["environment"].(string) != "prod" {
		c.Next()
		return
	}

	traceToken := c.GetHeader("X-Cloud-Trace-Context")
	if traceToken != "" {
		// TODO do something with the trace context
		c.Set("trace.id", traceToken)
	}

	protoHeader := c.GetHeader("X-Forwarded-Proto")
	if protoHeader == "" {
		Respond(c, http.StatusUnauthorized, "missing protocol header")
		return
	}
	if protoHeader != "https" {
		Respond(c, http.StatusUnauthorized, "refusing to serve unencrypted traffic")
		return
	}

	gatewayHeader := c.GetHeader("Authorization")
	if gatewayHeader == "" {
		Respond(c, http.StatusUnauthorized, "missing gateway authorization header")
		return
	}
	gatewayToken := strings.Split(gatewayHeader, " ")[1]
	clientHTTP := Global["client.http"].(*http.Client)
	response, err := clientHTTP.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", gatewayToken))
	if err != nil {
		Respond(c, http.StatusUnauthorized, "failed to reach gateway verification endpoint")
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "failed to read verification response")
		return
	}
	var responseMap map[string]interface{}
	json.Unmarshal(body, &responseMap)
	if _, ok := responseMap["error"]; ok {
		Respond(c, http.StatusUnauthorized, "failed to verify gateway token")
		return
	}
	if responseMap["iss"] != "https://accounts.google.com" {
		Respond(c, http.StatusUnauthorized, "failed to verify gateway token issuer")
		return
	}
	gatewayIdentity := Global["gateway.service_account"].(string)
	if responseMap["email"] != gatewayIdentity {
		Respond(c, http.StatusUnauthorized, "failed to verify gateway token identity")
		return
	}

	firebaseHeader := c.GetHeader("X-Forwarded-Authorization")
	if firebaseHeader == "" {
		Respond(c, http.StatusUnauthorized, "missing authorization user authorization header")
		return
	}
	firebaseToken := strings.Split(firebaseHeader, " ")[1]
	clientFirebase := Global["client.firebase"].(*auth.Client)
	token, err := clientFirebase.VerifyIDToken(ctx, firebaseToken)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "unable to verify user token")
		return
	}

	encoded := c.Request.Header.Get("X-Endpoint-API-UserInfo")
	if encoded == "" {
		Respond(c, http.StatusUnauthorized, "missing gateway user info header")
		return
	}

	bytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "failed to decode user info header")
		return
	}

	var caller Caller
	err = json.Unmarshal(bytes, &caller)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "failed to deserialize user info header")
		return
	}

	if token.UID != caller.ID {
		Respond(c, http.StatusUnauthorized, "mismatching inbound caller identities")
		return
	}

	// Verification OK

	c.Set("caller.email", caller.Email)
	c.Set("caller.id", caller.ID)
	c.Set("caller.name", caller.Name)
	c.Set("caller.picture", caller.PictureURL)
	c.Next()
}
