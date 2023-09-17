package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
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
			if instance.GetName() == config.InstanceName {
				return instance, nil
			}
		}
	}
	return nil, errors.New("no instance found")
}

func deployInstance(region string) (*computepb.Instance, error) {
	zone, err := findZone(region)
	if err != nil {
		return nil, fmt.Errorf("unable to find zone: %v", err)
	}

	disk, err := createDisk(zone)
	if err != nil {
		return nil, fmt.Errorf("failed to create disk: %v", err)
	}

	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client: %v", err)
	}
	defer client.Close()

	boot := true
	autoDelete := false
	machineType := fmt.Sprintf("zones/%s/machineTypes/%s", zone, config.MachineType)

	req := &computepb.InsertInstanceRequest{
		Project: config.Project,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name:            &config.InstanceName,
			ServiceAccounts: nil,
			MachineType:     &machineType,
			Disks: []*computepb.AttachedDisk{
				{
					Source:     disk.SelfLink,
					Boot:       &boot,
					AutoDelete: &autoDelete,
				},
			},
			Scheduling: &computepb.Scheduling{
				Preemptible:       &config.Preemtibility,
				ProvisioningModel: &config.ProvisioningModel,
				OnHostMaintenance: &config.OnHostMaintenance,
			},
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

	instance, err := getInstance()
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func destroyInstance(instance *computepb.Instance) error {
	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %v", err)
	}
	defer client.Close()

	req := &computepb.DeleteInstanceRequest{
		Project:  config.Project,
		Zone:     *instance.Zone,
		Instance: *instance.Name,
	}

	op, err := client.Delete(ctx, req)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	log.Printf("deleted instance: %s, in zone %s", *instance.Name, *instance.Zone)

	disk := *instance.Disks[0].Source

	snap, err := createSnapshot(disk)
	if err != nil {
		return err
	}

	log.Printf("created snapshot: %v", snap.Id)

	err = deleteDisk(disk, *instance.Zone)
	if err != nil {
		return err
	}

	log.Printf("deleted disk: %s, in zone %s", disk, *instance.Zone)

	return nil
}
