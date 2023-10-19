package main

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

func getSnapshot() (*computepb.Snapshot, error) {
	ctx := context.Background()
	client, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &computepb.GetSnapshotRequest{
		Project:  config.Project,
		Snapshot: config.SnapshotName,
	}

	snapshot, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func createSnapshot(diskName string) (*computepb.Snapshot, error) {
	ctx := context.Background()
	client, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &computepb.InsertSnapshotRequest{
		Project: config.Project,
		SnapshotResource: &computepb.Snapshot{
			Name:       &config.SnapshotName,
			SourceDisk: &diskName,
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

	snapshot, err := getSnapshot()
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func deleteSnapshot() error {
	ctx := context.Background()
	client, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.DeleteSnapshotRequest{
		Project:  config.Project,
		Snapshot: config.SnapshotName,
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
