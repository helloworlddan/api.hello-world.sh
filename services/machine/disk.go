package main

import (
	"context"
	"fmt"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

func createDisk(zone string) (*computepb.Disk, error) {
	name := "top"

	var diskSize int64
	diskSize = 20

	ctx := context.Background()
	client, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	snapshot := fmt.Sprintf(
		"projects/%s/global/snapshots/%s",
		config.Project,
		config.SnapshotName,
	)

	diskType := fmt.Sprintf(
		"projects/%s/zones/%s/diskTypes/%s",
		config.Project,
		zone,
		config.DiskType,
	)

	req := &computepb.InsertDiskRequest{
		Project: config.Project,
		Zone:    zone,
		DiskResource: &computepb.Disk{
			Name:           &name,
			SourceSnapshot: &snapshot,
			Type:           &diskType,
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

func deleteDisk(disk string, zone string) error {
	ctx := context.Background()
	client, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	zoneComponents := strings.Split(zone, "/")
	zone = zoneComponents[len(zoneComponents)-1]

	diskComponents := strings.Split(disk, "/")
	disk = diskComponents[len(diskComponents)-1]

	req := &computepb.DeleteDiskRequest{
		Project: config.Project,
		Zone:    zone,
		Disk:    disk,
	}

	op, err := client.Delete(ctx, req)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
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
		Disk:    config.DiskName,
		Project: config.Project,
	}
	disk, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return disk, nil
}
