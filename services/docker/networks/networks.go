package networks

import (
	"context"
	"medovukha/api/rest/v1/types"

	dc "medovukha/services/docker"

	"github.com/docker/docker/api/types/network"
)

func GetNetworkList(cli dc.IDockerClient) ([]types.NetworkBaseInfo, error) {
	ctx := context.Background()

	networks, err := cli.NetworkList(ctx, network.ListOptions{})
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
