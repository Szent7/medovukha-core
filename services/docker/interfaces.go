package docker

import (
	"context"
	"io"

	ts "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type IDockerClient interface {
	//Containers
	ContainerList(ctx context.Context, options container.ListOptions) ([]ts.Container, error)

	ContainerPause(ctx context.Context, containerID string) error

	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig,
		platform *ocispec.Platform, containerName string) (container.CreateResponse, error)

	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error

	ContainerUnpause(ctx context.Context, containerID string) error

	ContainerKill(ctx context.Context, containerID, signal string) error

	ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error

	ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error

	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error

	//Images
	ImagePull(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error)

	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)

	ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)

	ImageBuild(ctx context.Context, buildContext io.Reader, options ts.ImageBuildOptions) (ts.ImageBuildResponse, error)

	//Networks
	NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error)

	//Volumes
	VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)

	//Events
	Events(ctx context.Context, options events.ListOptions) (<-chan events.Message, <-chan error)
}
