package main

import (
	"context"
	"log"
	"net/http"
	"os"

	firestore "cloud.google.com/go/firestore"
	stackdriver "contrib.go.opencensus.io/exporter/stackdriver"
	firebase "firebase.google.com/go/v4"
	auth "firebase.google.com/go/v4/auth"
	gin "github.com/gin-gonic/gin"
	ochttp "go.opencensus.io/plugin/ochttp"
	propagation "go.opencensus.io/plugin/ochttp/propagation/b3"
	trace "go.opencensus.io/trace"
)

// Global map for shared resources
var Global map[string]interface{}

func main() {
	ctx := context.Background()
	configure(ctx)

	exporter := createTraceExporter()
	defer exporter.StopMetricsExporter()
	ctx, span := trace.StartSpan(ctx, "things.main")
	defer span.End()

	router := gin.Default()
	authorized := router.Group("/")
	authorized.Use(Authenticate)
	{
		authorized.GET("/things/:key", GetHandler)
		authorized.GET("/things/", ListHandler)
		authorized.POST("/things/", PostHandler)
		authorized.DELETE("/things/:key", DeleteHandler)
		authorized.PATCH("/things/:key", PatchHandler)
	}
	router.Run()
}

func configure(ctx context.Context) {
	Global = make(map[string]interface{})
	Global["environment"] = os.Getenv("HWSH_ENVIRONMENT")
	if Global["environment"].(string) == "" {
		Global["environment"] = "dev"
	}
	if Global["environment"].(string) == "prod" {
		gin.SetMode(gin.ReleaseMode)
		Global["gateway.service_account"] = os.Getenv("HWSH_GATEWAY_SA")
		if Global["gateway.service_account"].(string) == "" {
			log.Fatal("failed to read HWSH_GATEWAY_SA in production configuration")
		}
	}
	Global["project.id"] = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if Global["project.id"] == "" {
		log.Fatal("failed to read GOOGLE_CLOUD_PROJECT")
	}
	Global["client.http"] = createHTTPClient(ctx)
	Global["client.firebase"] = createFirebaseClient(ctx)
	Global["client.firestore"] = createFirestoreClient(ctx)
}

func createTraceExporter() *stackdriver.Exporter {
	projectID := Global["project.id"].(string)
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	exporter.StartMetricsExporter()
	return exporter
}

func createFirestoreClient(ctx context.Context) *firestore.Client {
	projectID := Global["project.id"].(string)
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	return client
}

func createFirebaseClient(ctx context.Context) *auth.Client {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("failed to create firebase app: %v", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("failed to create firebase client: %v", err)
	}
	return client
}

func createHTTPClient(ctx context.Context) *http.Client {
	client := &http.Client{
		Transport: &ochttp.Transport{
			Propagation: &propagation.HTTPFormat{},
		},
	}
	return client
}

// Respond terminates transaction with a standard error format
func Respond(c *gin.Context, code int, obj interface{}) {
	if code < 300 {
		if obj == nil {
			c.Status(code)
			c.Next()
			return
		}
		c.JSON(code, obj)
		c.Next()
		return
	}
	if obj == nil {
		c.Status(code)
		c.Abort()
		return
	}
	c.JSON(code, gin.H{"error": obj})
	c.Abort()
	return
}
