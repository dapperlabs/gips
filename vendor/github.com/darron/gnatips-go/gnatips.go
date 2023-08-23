package gnatips

import (
	"context"
	"fmt"
	"path"

	"google.golang.org/api/compute/v1"
)

func Get(ctx context.Context, project, region string) ([]string, error) {
	var ips []string

	// Create a Compute Engine service client
	client, err := compute.NewService(ctx)
	if err != nil {
		return ips, fmt.Errorf("Failed to create Compute Engine client: %w", err)
	}

	// List NAT gateways in the specified region
	natGateways, err := client.Routers.List(project, region).Do()
	if err != nil {
		return ips, fmt.Errorf("Failed to list NAT gateways: %w", err)
	}

	// Iterate through NAT gateways and gather external IP addresses
	for _, natGateway := range natGateways.Items {
		for _, nat := range natGateway.Nats {
			for _, natIP := range nat.NatIps {
				// Let's pull the last part of the URL to get the name of the address
				addressName := path.Base(natIP)
				address, err := client.Addresses.Get(project, region, addressName).Do()
				if err != nil {
					return ips, fmt.Errorf("Failed to get address information: %w", err)
				}
				ips = append(ips, address.Address)
			}
		}
	}
	return ips, nil
}
