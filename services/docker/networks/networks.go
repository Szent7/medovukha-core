package networks

import (
	"context"

	"github.com/Szent7/medovukha-core/ipc/types"

	dc "github.com/Szent7/medovukha-core/services/docker"

	"github.com/docker/docker/api/types/network"
)

func GetNetworkList(ctx context.Context, cli dc.IDockerClient) ([]types.NetworkBaseInfo, error) {
	networks, err := GetNetworkRawList(ctx, cli)
	if err != nil {
		return nil, err
	}

	netList := make([]types.NetworkBaseInfo, len(networks))
	for i, network := range networks {
		netList[i] = types.NetworkBaseInfo{
			Name:       network.Name,
			Id:         network.ID,
			Driver:     network.Driver,
			EnableIPv6: network.EnableIPv6,
			IPAMDriver: network.IPAM.Driver,
		}
		netList[i].Subnet = make([]string, len(network.IPAM.Config))
		netList[i].Gateway = make([]string, len(network.IPAM.Config))
		for j, netconfig := range network.IPAM.Config {
			netList[i].Subnet[j] = netconfig.Subnet
			netList[i].Gateway[j] = netconfig.Gateway
		}
		if network.Name == "none" || network.Name == "host" || network.Name == "bridge" {
			netList[i].DockerNetwork = true
		} else {
			netList[i].DockerNetwork = false
		}
	}

	return netList, nil
}

func GetNetwork(ctx context.Context, cli dc.IDockerClient, networkID string) (types.NetworkBaseInfo, error) {
	networkInspect, err := GetNetworkSummary(ctx, cli, networkID)
	if err != nil {
		return types.NetworkBaseInfo{}, err
	}

	ntwSummary := types.NetworkBaseInfo{
		Name:       networkInspect.Name,
		Id:         networkInspect.ID,
		Driver:     networkInspect.Driver,
		EnableIPv6: networkInspect.EnableIPv6,
		IPAMDriver: networkInspect.IPAM.Driver,
	}

	ntwSummary.Subnet = make([]string, len(networkInspect.IPAM.Config))
	ntwSummary.Gateway = make([]string, len(networkInspect.IPAM.Config))
	for i, netconfig := range networkInspect.IPAM.Config {
		ntwSummary.Subnet[i] = netconfig.Subnet
		ntwSummary.Gateway[i] = netconfig.Gateway
	}
	if networkInspect.Name == "none" || networkInspect.Name == "host" || networkInspect.Name == "bridge" {
		ntwSummary.DockerNetwork = true
	} else {
		ntwSummary.DockerNetwork = false
	}

	return ntwSummary, nil
}

func IsNetworkUsed(ctx context.Context, cli dc.IDockerClient, networkID string) (bool, error) {
	network, err := GetNetworkSummary(ctx, cli, networkID)
	if err != nil {
		return false, err
	}

	return len(network.Containers) != 0, nil
}

func GetNetworkRawList(ctx context.Context, cli dc.IDockerClient) ([]network.Summary, error) {
	return cli.NetworkList(ctx, network.ListOptions{})
}

func GetNetworkSummary(ctx context.Context, cli dc.IDockerClient, networkID string) (network.Summary, error) {
	return cli.NetworkInspect(ctx, networkID, network.InspectOptions{Verbose: true})
}

func RemoveNetworkByID(ctx context.Context, cli dc.IDockerClient, networkID string) error {
	return cli.NetworkRemove(ctx, networkID)
}
