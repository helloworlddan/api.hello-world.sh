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
	haver "github.com/umahmood/haversine"
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

		current := strings.Split(resp.GetRegion(), "/")

		if current[len(current)-1] != region {
			continue
		}
		return resp.GetName(), nil
	}

	return "", fmt.Errorf("no suitable zone found for region: %s\n", region)
}

func ipSource(r *http.Request) string {
	addr := r.Header.Get("TOP-DEBUG-IP")
	if addr == "" {
		addr = strings.Split(r.RemoteAddr, ":")[0]
	}
	return addr
}

func ipLocation(ip string) (haver.Coord, error) {
	var lat float64
	var lon float64

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://ipinfo.io/%s", ip), nil)
	if err != nil {
		return haver.Coord{}, fmt.Errorf("failed to construct request: %v\n", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return haver.Coord{}, fmt.Errorf("failed to resolve ip address: %v\n", err)
	}
	defer res.Body.Close()

	var target map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&target)
	if err != nil {
		return haver.Coord{}, fmt.Errorf("failed to decode JSON response: %v\n", err)
	}

	location := target["loc"].(string)

	locationComponents := strings.Split(location, ",")

	lat, err = strconv.ParseFloat(locationComponents[0], 64)
	if err != nil {
		return haver.Coord{}, fmt.Errorf("failed to parse latitude from response: %v\n", err)
	}

	lon, err = strconv.ParseFloat(locationComponents[1], 64)
	if err != nil {
		return haver.Coord{}, fmt.Errorf("failed to parse longitude from response: %v\n", err)
	}

	return haver.Coord{Lat: lat, Lon: lon}, nil
}

func seedRegions() map[string]haver.Coord {
	regions := make(map[string]haver.Coord)

	regions["europe-central2"] = haver.Coord{Lat: 52.229, Lon: 21.012}          // Warsaw
	regions["europe-north1"] = haver.Coord{Lat: 60.569, Lon: 27.187}            // Hamina
	regions["europe-west1"] = haver.Coord{Lat: 50.470, Lon: 3.817}              // St. Ghislain
	regions["europe-west2"] = haver.Coord{Lat: 51.507, Lon: 0.127}              // London
	regions["europe-west3"] = haver.Coord{Lat: 50.110, Lon: 8.682}              // Frankfurt
	regions["europe-west4"] = haver.Coord{Lat: 53.438, Lon: 6.835}              // Eemshaven
	regions["europe-west6"] = haver.Coord{Lat: 47.376, Lon: 8.541}              // Zurich
	regions["europe-west8"] = haver.Coord{Lat: 45.464, Lon: 9.190}              // Milan
	regions["europe-west9"] = haver.Coord{Lat: 48.856, Lon: 2.352}              // Paris
	regions["europe-west10"] = haver.Coord{Lat: 52.520, Lon: 13.404}            // Berlin
	regions["europe-west12"] = haver.Coord{Lat: 47.070, Lon: 7.686}             // Turin
	regions["europe-southwest1"] = haver.Coord{Lat: 40.416, Lon: 3.703}         // Madrid
	regions["us-central1"] = haver.Coord{Lat: 41.261, Lon: -95.860}             // Council Bluffs
	regions["us-west1"] = haver.Coord{Lat: 43.804, Lon: -120.554}               // Oregon
	regions["us-west2"] = haver.Coord{Lat: 34.054, Lon: -118.242}               // Los Angeles
	regions["us-west3"] = haver.Coord{Lat: 40.760, Lon: -111.891}               // Salt Lake City
	regions["us-west4"] = haver.Coord{Lat: 36.171, Lon: -115.139}               // Las Vegas
	regions["us-east1"] = haver.Coord{Lat: 33.126, Lon: -80.008}                // Berkely County
	regions["us-east4"] = haver.Coord{Lat: 39.076, Lon: -77.653}                // Loudoun County
	regions["us-east5"] = haver.Coord{Lat: 39.961, Lon: -82.998}                // Columbus
	regions["us-south1"] = haver.Coord{Lat: 32.776, Lon: -96.797}               // Dallas
	regions["northamerica-northeast1"] = haver.Coord{Lat: 45.501, Lon: -73.567} // Montreal
	regions["northamerica-northeast2"] = haver.Coord{Lat: 43.653, Lon: -79.383} // Toronto
	regions["southamerica-west1"] = haver.Coord{Lat: -33.357, Lon: -70.729}     // Quilicura
	regions["southamerica-east1"] = haver.Coord{Lat: -23.555, Lon: -46.639}     // Sao Paolo
	regions["asia-south1"] = haver.Coord{Lat: 18.975, Lon: 72.825}              // Mumbai
	regions["asia-south2"] = haver.Coord{Lat: 28.684, Lon: 77.222}              // Delhi
	regions["asia-southeast1"] = haver.Coord{Lat: 1.366, Lon: 103.800}          // Singapore
	regions["asia-southeast2"] = haver.Coord{Lat: -6.174, Lon: 106.829}         // Jakarta
	regions["asia-east1"] = haver.Coord{Lat: 24.066, Lon: 120.533}              // Changhua County
	regions["asia-east2"] = haver.Coord{Lat: 22.319, Lon: 114.169}              // Hong Kong
	regions["asia-northeast1"] = haver.Coord{Lat: 35.685, Lon: 139.751}         // Tokio
	regions["asia-northeast2"] = haver.Coord{Lat: 34.666, Lon: 135.500}         // Osaka
	regions["asia-northeast3"] = haver.Coord{Lat: 37.565, Lon: 126.565}         // Seoul
	regions["australia-southeast1"] = haver.Coord{Lat: -33.869, Lon: 151.209}   // Sydney
	regions["australia-southeast2"] = haver.Coord{Lat: -37.813, Lon: 144.963}   // Melbourne

	return regions
}

func closestRegion(source haver.Coord) (string, float64) {
	var distance float64
	distance = 1000000

	var closest string

	for name, coordinates := range config.Regions {
		_, d := haver.Distance(source, coordinates)

		if d >= distance {
			continue
		}

		distance = d
		closest = name
	}

	return closest, distance
}
