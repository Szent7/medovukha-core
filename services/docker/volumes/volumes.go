package volumes

import (
	"context"

	"github.com/Szent7/medovukha-core/ipc/types"

	dc "github.com/Szent7/medovukha-core/services/docker"

	"github.com/docker/docker/api/types/volume"
)

func GetVolumeList(cli dc.IDockerClient) ([]types.VolumeBaseInfo, error) {
	ctx := context.Background()

	volumes, err := GetVolumeRawList(ctx, cli)
	if err != nil {
		return nil, err
	}

	volList := make([]types.VolumeBaseInfo, len(volumes.Volumes))
	for i, volume := range volumes.Volumes {
		volList[i] = types.VolumeBaseInfo{
			Name:       volume.Name,
			Driver:     volume.Driver,
			Mountpoint: volume.Mountpoint,
			Created:    volume.CreatedAt,
		}
	}

	return volList, nil
}

func GetVolume(ctx context.Context, cli dc.IDockerClient, volumeID string) (types.VolumeBaseInfo, error) {
	volumeInspect, err := cli.VolumeInspect(ctx, volumeID)
	if err != nil {
		return types.VolumeBaseInfo{}, err
	}

	volSummary := types.VolumeBaseInfo{
		Name:       volumeInspect.Name,
		Driver:     volumeInspect.Driver,
		Mountpoint: volumeInspect.Mountpoint,
		Created:    volumeInspect.CreatedAt,
	}

	return volSummary, nil
}

func GetVolumeRawList(ctx context.Context, cli dc.IDockerClient) (volume.ListResponse, error) {
	return cli.VolumeList(ctx, volume.ListOptions{})
}

func RemoveVolumeByID(ctx context.Context, cli dc.IDockerClient, volumeID string, force bool) error {
	return cli.VolumeRemove(ctx, volumeID, force)
}
