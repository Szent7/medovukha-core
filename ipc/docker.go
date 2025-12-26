package ipc

import (
	"context"
	"log"

	"github.com/docker/docker/client"
	"google.golang.org/grpc"

	"errors"
	"fmt"

	dockerpb "github.com/Szent7/medovukha-core/api/docker/v1"
	"github.com/Szent7/medovukha-core/ipc/types"

	"net/url"
	"os"
	"path/filepath"
	"strings"

	containers "github.com/Szent7/medovukha-core/services/docker/containers"
	images "github.com/Szent7/medovukha-core/services/docker/images"
	networks "github.com/Szent7/medovukha-core/services/docker/networks"
	volumes "github.com/Szent7/medovukha-core/services/docker/volumes"
	git "github.com/Szent7/medovukha-core/services/git"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

type DockerService struct {
	dockerpb.UnimplementedDockerServiceServer
	cli *client.Client
}

func NewDockerCore() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &DockerService{cli: cli}, nil
}

func (d *DockerService) Close() error {
	return d.cli.Close()
}

// Containers
func (d *DockerService) GetContainerList(ctx context.Context, req *dockerpb.Empty) (*dockerpb.ListContainerBaseInfo, error) {
	containerList, err := containers.GetContainerBaseInfoList(d.cli)
	if err != nil {
		fmt.Printf("GetContainerList error: %s\n", err.Error())
		return nil, err
	}

	var serializedList dockerpb.ListContainerBaseInfo
	serializedList.Items = make([]*dockerpb.ContainerBaseInfo, len(containerList))
	for i := range containerList {
		serializedList.Items[i] = &dockerpb.ContainerBaseInfo{
			Id:        containerList[i].Id,
			Names:     containerList[i].Names,
			ImageName: containerList[i].ImageName,
			Ports:     convertPorts(containerList[i].Ports),
			Created:   containerList[i].Created,
			State:     containerList[i].State,
		}
	}

	return &serializedList, nil
}

func (d *DockerService) PauseContainerByID(ctx context.Context, req *dockerpb.BaseID) (*dockerpb.BaseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("PauseContainerByID error: nil request")
	}

	if err := containers.PauseContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("PauseContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("PauseContainerByID error: %s", err.Error())
	}

	return &dockerpb.BaseResponse{Message: "Paused: " + req.Id}, nil
}

func (d *DockerService) UnpauseContainerByID(ctx context.Context, req *dockerpb.BaseID) (*dockerpb.BaseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("UnpauseContainerByID error: nil request")
	}

	if err := containers.UnpauseContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("UnpauseContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("UnpauseContainerByID error: %s", err.Error())
	}

	return &dockerpb.BaseResponse{Message: "Unpaused: " + req.Id}, nil
}

func (d *DockerService) KillContainerByID(ctx context.Context, req *dockerpb.BaseID) (*dockerpb.BaseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("KillContainerByID error: nil request")
	}

	if err := containers.KillContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("KillContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("KillContainerByID error: %s", err.Error())
	}

	return &dockerpb.BaseResponse{Message: "Killed: " + req.Id}, nil
}

func (d *DockerService) StartContainerByID(ctx context.Context, req *dockerpb.BaseID) (*dockerpb.BaseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("StartContainerByID error: nil request")
	}

	if err := containers.StartContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("StartContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("StartContainerByID error: %s", err.Error())
	}

	return &dockerpb.BaseResponse{Message: "Started: " + req.Id}, nil
}

func (d *DockerService) StopContainerByID(ctx context.Context, req *dockerpb.BaseID) (*dockerpb.BaseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("StopContainerByID error: nil request")
	}

	if err := containers.StopContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("StopContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("StopContainerByID error: %s", err.Error())
	}

	return &dockerpb.BaseResponse{Message: "Stopped: " + req.Id}, nil
}

func (d *DockerService) RestartContainerByID(ctx context.Context, req *dockerpb.BaseID) (*dockerpb.BaseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("RestartContainerByID error: nil request")
	}

	if err := containers.RestartContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("RestartContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("RestartContainerByID error: %s", err.Error())
	}

	return &dockerpb.BaseResponse{Message: "Restarted: " + req.Id}, nil
}

func (d *DockerService) RemoveContainerByID(ctx context.Context, req *dockerpb.BaseID) (*dockerpb.BaseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("RemoveContainerByID error: nil request")
	}

	if err := containers.RemoveContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("RemoveContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("RemoveContainerByID error: %s", err.Error())
	}

	return &dockerpb.BaseResponse{Message: "Removed: " + req.Id}, nil
}

// Images
func (d *DockerService) GetImageList(ctx context.Context, req *dockerpb.Empty) (*dockerpb.ListImageBaseInfo, error) {
	imageList, err := images.GetImageList(d.cli)
	if err != nil {
		fmt.Printf("GetImageList error: %s\n", err.Error())
		return nil, fmt.Errorf("GetImageList error: %s", err.Error())
	}

	var serializedList dockerpb.ListImageBaseInfo
	serializedList.Items = make([]*dockerpb.ImageBaseInfo, len(imageList))
	for i := range imageList {
		serializedList.Items[i] = &dockerpb.ImageBaseInfo{
			Id:      imageList[i].Id,
			Tags:    imageList[i].Tags,
			Size:    imageList[i].Size,
			Created: imageList[i].Created,
		}
	}

	return &serializedList, nil
}

// Networks
func (d *DockerService) GetNetworkList(ctx context.Context, req *dockerpb.Empty) (*dockerpb.ListNetworkBaseInfo, error) {
	networkList, err := networks.GetNetworkList(d.cli)
	if err != nil {
		fmt.Printf("GetNetworkList error: %s\n", err.Error())
		return nil, fmt.Errorf("GetNetworkList error: %s", err.Error())
	}

	var serializedList dockerpb.ListNetworkBaseInfo
	serializedList.Items = make([]*dockerpb.NetworkBaseInfo, len(networkList))
	for i := range networkList {
		serializedList.Items[i] = &dockerpb.NetworkBaseInfo{
			Name:          networkList[i].Name,
			Id:            networkList[i].Id,
			Driver:        networkList[i].Driver,
			EnableIpv6:    networkList[i].EnableIPv6,
			IpamDriver:    networkList[i].IPAMDriver,
			Subnet:        networkList[i].Subnet,
			Gateway:       networkList[i].Gateway,
			Attachable:    networkList[i].Attachable,
			DockerNetwork: networkList[i].DockerNetwork,
		}
	}

	return &serializedList, nil
}

// Volumes
func (d *DockerService) GetVolumeList(ctx context.Context, req *dockerpb.Empty) (*dockerpb.ListVolumeBaseInfo, error) {
	volumeList, err := volumes.GetVolumeList(d.cli)
	if err != nil {
		fmt.Printf("GetVolumeList error: %s\n", err.Error())
		return nil, fmt.Errorf("GetVolumeList error: %s", err.Error())
	}

	var serializedList dockerpb.ListVolumeBaseInfo
	serializedList.Items = make([]*dockerpb.VolumeBaseInfo, len(volumeList))
	for i := range volumeList {
		serializedList.Items[i] = &dockerpb.VolumeBaseInfo{
			Name:       volumeList[i].Name,
			Driver:     volumeList[i].Driver,
			Mountpoint: volumeList[i].Mountpoint,
			Created:    volumeList[i].Created,
		}
	}

	return &serializedList, nil
}

// Deploy
func (d *DockerService) CreateFromGit(ctx context.Context, req *dockerpb.DeployFromGit) (*dockerpb.BaseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("CreateFromGit error: nil request")
	}
	//
	//
	// Clonning repo from git into temp dir
	//
	//
	tempDir, err := os.MkdirTemp("", "docker-repo")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Created tempDir: " + tempDir)
	defer os.RemoveAll(tempDir)
	defer fmt.Println("Deleted tempDir: " + tempDir)

	if err := git.CloneRepo(&git.RepoCloner{}, req.Url, tempDir); err != nil {
		fmt.Printf("CloneRepo error: %s\n", err.Error())
		return nil, fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	//
	//
	// Check Dockerfile in request
	//
	//
	dockerfileDir := filepath.Join(tempDir, "Dockerfile")
	dockercomposeDir := filepath.Join(tempDir, "docker-compose.yml")
	if req.Dockerfile == "" {
		// Dockerfile is not specified in the request, check in the directory
		// If Dockerfile doesn`t exist, throw error
		if !fileExists(dockerfileDir) {
			fmt.Println("Dockerfile empty error")
			return nil, fmt.Errorf("CreateFromGit error: Dockerfile empty")
		}
	} else {
		// If Dockerfile specified in the request, rewrite/create new Dockerfile
		if err := createFile(dockerfileDir, []byte(req.Dockerfile)); err != nil {
			fmt.Printf("Dockerfile write error: %s\n", err.Error())
			return nil, fmt.Errorf("CreateFromGit error: %s", err.Error())
		} else {
			fmt.Println("Dockerfile redefined")
		}
	}
	//
	//
	// Parse tags from repo (for image name)
	//
	//
	repo, err := repoFromURL(req.Url)
	if err != nil {
		fmt.Printf("Repo error: %s\n", err.Error())
		return nil, fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	tags := []string{repo + ":latest"} //args.URL
	//
	//
	// Delete old containers
	//
	//
	if err := containers.RemoveContainerByImage(ctx, d.cli, tags[0]); err != nil {
		fmt.Printf("Error remove container: %s\n", err.Error())
		return nil, fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	//
	//
	// Delete old images
	//
	//
	if err := images.RemoveImageByTag(ctx, d.cli, tags[0]); err != nil {
		fmt.Printf("Error remove image: %s\n", err.Error())
		return nil, fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	//
	//
	// Build image
	//
	//
	newImageId, err := images.BuildImageNew(d.cli, tempDir, tags)
	if err != nil {
		fmt.Printf("Docker client error: %s\n", err.Error())
		return nil, fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	//
	//
	// Clear build cache
	//
	//
	pruneFilters := filters.NewArgs()
	pruneFilters.Add("dangling", "true")

	_, err = d.cli.ImagesPrune(ctx, pruneFilters)
	if err != nil {
		fmt.Printf("Docker client error: %s\n", err.Error())
		return nil, fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	_, err = d.cli.BuildCachePrune(ctx, build.CachePruneOptions{All: true})
	if err != nil {
		fmt.Printf("Docker client error: %s\n", err.Error())
		return nil, fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	//
	//
	// Check DockerCompose in request
	//
	//
	if req.DockerCompose == "" {
		// DockerCompose is not specified in the request, check DockerRun in the request
		if req.DockerRun == "" {
			// DockerRun is not specified in the request, check DockerCompose in the directory
			if !fileExists(dockercomposeDir) {
				fmt.Println("Created image, but not command to launch container")
				return &dockerpb.BaseResponse{Message: "created but not launched: " + newImageId}, nil
			} else {
				if err := containers.ExecDockerComposeUp(dockercomposeDir); err != nil {
					return &dockerpb.BaseResponse{Message: "the image was created, but an error occurred when starting the container (dockerCompose): " + newImageId}, nil
				}
			}
		} else {
			if err := containers.ExecDockerRun(req.DockerRun); err != nil {
				return &dockerpb.BaseResponse{Message: "the image was created, but an error occurred when starting the container (dockerRun): " + newImageId}, nil
			}
		}
	} else {
		// If Dockerfile specified in the request, rewrite/create new Dockerfile
		if err := createFile(dockercomposeDir, []byte(req.DockerCompose)); err != nil {
			fmt.Printf("DockerCompose write error: %s\n", err.Error())
			return nil, fmt.Errorf("CreateFromGit error: %s", err.Error())
		} else {
			fmt.Println("DockerCompose redefined")
			if err := containers.ExecDockerComposeUp(dockercomposeDir); err != nil {
				return &dockerpb.BaseResponse{Message: "the image was created, but an error occurred when starting the container (dockerCompose): " + newImageId}, nil
			}
		}
	}

	return &dockerpb.BaseResponse{Message: "created and launched: " + newImageId}, nil
}

// events
func (d *DockerService) GetContainerState(req *dockerpb.Empty, stream grpc.ServerStreamingServer[dockerpb.ContainerState]) error {
	ctx := stream.Context()

	eventCh := make(chan events.Message)
	errCh := make(chan error)
	go containers.EventStream(ctx, d.cli, eventCh, errCh)

	fmt.Println("Start gRPC Server-Side streaming")

	for {
		select {
		case <-ctx.Done():
			{
				fmt.Println("Stop gRPC Server-Side streaming")
				return nil
			}
		case err := <-errCh:
			{
				if err != nil {
					fmt.Println("Stop gRPC Server-Side streaming")
					return err
				}
			}
		case event := <-eventCh:
			{
				// fmt.Println("New event in GetContainerState")
				resp := &dockerpb.ContainerState{
					Type:   string(event.Type),
					Action: string(event.Action),
					Actor: &dockerpb.Actor{
						Id:         event.Actor.ID,
						Attributes: event.Actor.Attributes,
					},
				}
				if err := stream.Send(resp); err != nil {
					log.Printf("Send error: %s", err.Error())
					fmt.Println("Stop gRPC Server-Side streaming")
					return err
				}
			}
		}
	}
}

// common
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	fmt.Println("ошибка при проверке:", err)
	return false
}

func createFile(path string, content []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		return err
	}

	return nil
}

func repoFromURL(raw string) (repo string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	// u.Path выглядит как "/owner/repo" (или "/owner/repo/…")
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid GitHub URL: %s", raw)
	}
	return parts[0] + "/" + parts[1], nil
}

func convertPorts(jsonPorts []types.Port) []*dockerpb.Port {
	if len(jsonPorts) == 0 {
		return nil
	}

	grpcPorts := make([]*dockerpb.Port, len(jsonPorts))
	for i := range jsonPorts {
		grpcPorts[i] = &dockerpb.Port{
			Ip:          jsonPorts[i].IP,
			PrivatePort: uint32(jsonPorts[i].PrivatePort),
			PublicPort:  uint32(jsonPorts[i].PublicPort),
			Type:        jsonPorts[i].Type,
		}
	}
	return grpcPorts
}
