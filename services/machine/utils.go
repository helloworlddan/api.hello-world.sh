package main

import (
	"context"
	"encoding/json"
	"errors"
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

func getInstance() (*computepb.Instance, error) {
	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client: %v", err)
	}
	defer client.Close()

	req := &computepb.AggregatedListInstancesRequest{
		Project:    config.Project,
		MaxResults: proto.Uint32(3),
	}

	it := client.AggregatedList(ctx, req)
	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read instances: %v/n", err)
		}
		instances := pair.Value.Instances
		for _, instance := range instances {
			if instance.GetName() == config.Machine {
				return instance, nil
			}
		}
	}
	return nil, errors.New("no instance found")
}

func deployInstance(region string) error {
	zone, err := findZone(region)
	if err != nil {
		return fmt.Errorf("unable to find zone: %v", err)
	}

	disk, err := createDisk(zone)
	if err != nil {
		return fmt.Errorf("failed to create disk: %v", err)
	}

	// TODO continue to create instance
	_ = disk

	return nil
}

func destroyInstance(instance *computepb.Instance) error {
	return nil
}

func createDisk(zone string) (*computepb.Disk, error) {
	var diskSize int64
	diskSize = 20

	ctx := context.Background()
	client, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	snapshot := fmt.Sprintf("/projects/%s/global/snapshots/%s", config.Project, config.Snapshot)

	req := &computepb.InsertDiskRequest{
		Project: config.Project,
		Zone:    zone,
		DiskResource: &computepb.Disk{
			SourceSnapshot: &snapshot,
			SizeGb:         &diskSize,
		},
	}
	op, err := client.Insert(ctx, req)
	if err != nil {
		return nil, err
	}

	err = op.Wait(ctx)
	if err != nil {
		return nil, err
	}

	disk, err := getDisk(zone)
	if err != nil {
		return nil, err
	}

	return disk, nil
}

func getDisk(zone string) (*computepb.Disk, error) {
	ctx := context.Background()
	client, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &computepb.GetDiskRequest{
		Zone:    zone,
		Disk:    config.Disk,
		Project: config.Project,
	}
	disk, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return disk, nil
}

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
	regions["europe-west10"] = haversine.Coord{Lat: 53.0, Lon: 9.0}

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
