package main

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

func snapshotName() string {
	return fmt.Sprintf("/projects/%s/global/snapshots/%s", config.Project, config.SnapshotName)
}

func getSnapshot() (*computepb.Snapshot, error) {
	ctx := context.Background()
	client, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &computepb.GetSnapshotRequest{
		Project:  config.Project,
		Snapshot: snapshotName(),
	}

	snapshot, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func createSnapshot(disk *computepb.Disk) (*computepb.Snapshot, error) {
	// TODO continue snapshot creation
	return nil, nil
}
