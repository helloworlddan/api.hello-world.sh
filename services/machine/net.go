package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	haversine "github.com/umahmood/haversine"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
)

func findZone(region string) (string, error) {
	ctx := context.Background()
	client, err := compute.NewZonesRESTClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	req := &computepb.ListZonesRequest{
		Project:    config.Project,
		MaxResults: proto.Uint32(3),
	}
	it := client.List(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}

		if resp.GetRegion() != region {
			continue
		}
		return resp.GetName(), nil
	}
	return "", fmt.Errorf("no suitable zone found for region: %s\n", region)
}

func ipSource(r *http.Request) string {
	addr := r.Header.Get("X-Real-Ip")
	if addr == "" {
		addr = r.Header.Get("X-Forwarded-For")
	}
	if addr == "" {
		addr = r.RemoteAddr
	}
	return addr
}

func ipLocation(ip string) (haversine.Coord, error) {
	var lat float64
	var lon float64

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://ipinfo.io/%s", ip), nil)
	if err != nil {
		return haversine.Coord{}, fmt.Errorf("failed to construct request: %v\n", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return haversine.Coord{}, fmt.Errorf("failed to resolve ip address: %v\n", err)
	}
	defer res.Body.Close()

	var target map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&target)
	if err != nil {
		return haversine.Coord{}, fmt.Errorf("failed to decode JSON response: %v\n", err)
	}

	location := target["loc"].(string)

	locationComponents := strings.Split(location, ",")

	lat, err = strconv.ParseFloat(locationComponents[0], 64)
	if err != nil {
		return haversine.Coord{}, fmt.Errorf("failed to parse latitude from response: %v\n", err)
	}

	lon, err = strconv.ParseFloat(locationComponents[1], 64)
	if err != nil {
		return haversine.Coord{}, fmt.Errorf("failed to parse longitude from response: %v\n", err)
	}

	return haversine.Coord{Lat: lat, Lon: lon}, nil
}

func seedRegions() map[string]haversine.Coord {
	regions := make(map[string]haversine.Coord)

	regions["europe-central1"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-north1"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-west1"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-west2"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-west3"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-west4"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-west5"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-west6"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-west7"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-west8"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-west9"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-west10"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["europe-southwest1"] = haversine.Coord{Lat: 0.0, Lon: 0.0}

	regions["us-central1"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["us-west1"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["us-west2"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["us-east1"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["us-east2"] = haversine.Coord{Lat: 0.0, Lon: 0.0}

	regions["asia-southwest1"] = haversine.Coord{Lat: 0.0, Lon: 0.0}
	regions["asia-southwest2"] = haversine.Coord{Lat: 0.0, Lon: 0.0}

	return regions
}

func closestRegion(source haversine.Coord) (string, float64) {
	var distance float64
	distance = 1000000

	var closest string

	for name, coordinates := range config.Regions {
		_, d := haversine.Distance(source, coordinates)

		if d >= distance {
			continue
		}

		distance = d
		closest = name
	}

	return closest, distance
}
