package ipc

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"google.golang.org/grpc"

	dockerpb "github.com/Szent7/medovukha-core/proto/docker/v1"

	"os"
	"path/filepath"

	"github.com/Szent7/medovukha-core/services/common"
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

	logChans map[string]chan string //key - BuildID; value - log
	logMu    sync.RWMutex
}

func NewDockerCore() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &DockerService{
		cli:      cli,
		logChans: make(map[string]chan string),
	}, nil
}

func (d *DockerService) Close() error {
	return d.cli.Close()
}

// Containers
func (d *DockerService) GetContainerList(ctx context.Context, req *dockerpb.GetContainerListRequest) (*dockerpb.GetContainerListResponse, error) {
	containerList, err := containers.GetContainerBaseInfoList(ctx, d.cli)
	if err != nil {
		fmt.Printf("GetContainerList error: %s\n", err.Error())
		return nil, err
	}

	var serializedList dockerpb.GetContainerListResponse
	serializedList.Items = make([]*dockerpb.ContainerBaseInfo, len(containerList))
	for i := range containerList {
		serializedList.Items[i] = &dockerpb.ContainerBaseInfo{
			Id:        containerList[i].Id,
			Names:     containerList[i].Names,
			ImageName: containerList[i].ImageName,
			Ports:     common.ConvertPorts(containerList[i].Ports),
			Created:   containerList[i].Created,
			State:     containerList[i].State,
		}
	}

	return &serializedList, nil
}

func (d *DockerService) GetContainerByID(ctx context.Context, req *dockerpb.GetContainerByIDRequest) (*dockerpb.GetContainerByIDResponse, error) {
	container, err := containers.GetContainer(ctx, d.cli, req.Id)
	if err != nil {
		fmt.Printf("GetContainer error: %s\n", err.Error())
		return nil, err
	}

	serializedContainer := dockerpb.GetContainerByIDResponse{
		Item: &dockerpb.ContainerBaseInfo{
			Id:        container.Id,
			Names:     container.Names,
			ImageName: container.ImageName,
			Ports:     common.ConvertPorts(container.Ports),
			Created:   container.Created,
			State:     container.State,
		},
	}

	return &serializedContainer, nil
}

func (d *DockerService) PauseContainerByID(ctx context.Context, req *dockerpb.PauseContainerByIDRequest) (*dockerpb.PauseContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("PauseContainerByID error: nil request")
	}

	if err := containers.PauseContainerByID(ctx, d.cli, req.Id); err != nil {
		log.Printf("PauseContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("PauseContainerByID error: %s", err.Error())
	}

	return &dockerpb.PauseContainerByIDResponse{Message: "Paused: " + req.Id}, nil
}

func (d *DockerService) UnpauseContainerByID(ctx context.Context, req *dockerpb.UnpauseContainerByIDRequest) (*dockerpb.UnpauseContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("UnpauseContainerByID error: nil request")
	}

	if err := containers.UnpauseContainerByID(ctx, d.cli, req.Id); err != nil {
		log.Printf("UnpauseContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("UnpauseContainerByID error: %s", err.Error())
	}

	return &dockerpb.UnpauseContainerByIDResponse{Message: "Unpaused: " + req.Id}, nil
}

func (d *DockerService) KillContainerByID(ctx context.Context, req *dockerpb.KillContainerByIDRequest) (*dockerpb.KillContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("KillContainerByID error: nil request")
	}

	if err := containers.KillContainerByID(ctx, d.cli, req.Id); err != nil {
		log.Printf("KillContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("KillContainerByID error: %s", err.Error())
	}

	return &dockerpb.KillContainerByIDResponse{Message: "Killed: " + req.Id}, nil
}

func (d *DockerService) StartContainerByID(ctx context.Context, req *dockerpb.StartContainerByIDRequest) (*dockerpb.StartContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("StartContainerByID error: nil request")
	}

	if err := containers.StartContainerByID(ctx, d.cli, req.Id); err != nil {
		log.Printf("StartContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("StartContainerByID error: %s", err.Error())
	}

	return &dockerpb.StartContainerByIDResponse{Message: "Started: " + req.Id}, nil
}

func (d *DockerService) StopContainerByID(ctx context.Context, req *dockerpb.StopContainerByIDRequest) (*dockerpb.StopContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("StopContainerByID error: nil request")
	}

	if err := containers.StopContainerByID(ctx, d.cli, req.Id); err != nil {
		log.Printf("StopContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("StopContainerByID error: %s", err.Error())
	}

	return &dockerpb.StopContainerByIDResponse{Message: "Stopped: " + req.Id}, nil
}

func (d *DockerService) RestartContainerByID(ctx context.Context, req *dockerpb.RestartContainerByIDRequest) (*dockerpb.RestartContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("RestartContainerByID error: nil request")
	}

	if err := containers.RestartContainerByID(ctx, d.cli, req.Id); err != nil {
		log.Printf("RestartContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("RestartContainerByID error: %s", err.Error())
	}

	return &dockerpb.RestartContainerByIDResponse{Message: "Restarted: " + req.Id}, nil
}

func (d *DockerService) RemoveContainerByID(ctx context.Context, req *dockerpb.RemoveContainerByIDRequest) (*dockerpb.RemoveContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("RemoveContainerByID error: nil request")
	}

	if err := containers.RemoveContainerByID(ctx, d.cli, req.Id); err != nil {
		log.Printf("RemoveContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("RemoveContainerByID error: %s", err.Error())
	}

	return &dockerpb.RemoveContainerByIDResponse{Message: "Removed: " + req.Id}, nil
}

// Images
func (d *DockerService) GetImageList(ctx context.Context, req *dockerpb.GetImageListRequest) (*dockerpb.GetImageListResponse, error) {
	imageList, err := images.GetImageList(ctx, d.cli)
	if err != nil {
		log.Printf("GetImageList error: %s\n", err.Error())
		return nil, fmt.Errorf("GetImageList error: %s", err.Error())
	}

	var serializedList dockerpb.GetImageListResponse
	serializedList.Items = make([]*dockerpb.ImageBaseInfo, len(imageList))
	for i := range imageList {
		used, err := containers.IsImageUsed(ctx, d.cli, imageList[i].Id)
		if err != nil {
			log.Printf("GetImageList error: %s\n", err.Error())
			return nil, fmt.Errorf("GetImageList error: %s", err.Error())
		}
		serializedList.Items[i] = &dockerpb.ImageBaseInfo{
			Id:      imageList[i].Id,
			Tags:    imageList[i].Tags,
			Size:    imageList[i].Size,
			Created: imageList[i].Created,
			IsUsed:  used,
		}
	}

	return &serializedList, nil
}

func (d *DockerService) GetImageByID(ctx context.Context, req *dockerpb.GetImageByIDRequest) (*dockerpb.GetImageByIDResponse, error) {
	image, err := images.GetImage(ctx, d.cli, req.Id)
	if err != nil {
		fmt.Printf("GetImage error: %s\n", err.Error())
		return nil, err
	}

	used, err := containers.IsImageUsed(ctx, d.cli, req.Id)
	if err != nil {
		log.Printf("GetImageByID error: %s\n", err.Error())
		return nil, fmt.Errorf("GetImageByID error: %s", err.Error())
	}

	serializedImage := dockerpb.GetImageByIDResponse{
		Item: &dockerpb.ImageBaseInfo{
			Id:      image.Id,
			Tags:    image.Tags,
			Size:    image.Size,
			Created: image.Created,
			IsUsed:  used,
		},
	}

	return &serializedImage, nil
}

func (d *DockerService) RemoveImage(ctx context.Context, req *dockerpb.RemoveImageRequest) (*dockerpb.RemoveImageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("RemoveImage error: nil request")
	}

	if req.Force {
		if err := images.RemoveImageByID(ctx, d.cli, req.Id, req.Force); err != nil {
			log.Printf("RemoveImage error: %s\n", err.Error())
			return nil, fmt.Errorf("RemoveImage error: %s", err.Error())
		}
	} else {
		used, err := containers.IsImageUsed(ctx, d.cli, req.Id)
		if err != nil {
			log.Printf("RemoveImage error: %s\n", err.Error())
			return nil, fmt.Errorf("RemoveImage error: %s", err.Error())
		}

		if used {
			return nil, fmt.Errorf("image is in use")
		}

		if err := images.RemoveImageByID(ctx, d.cli, req.Id, false); err != nil {
			log.Printf("RemoveImage error: %s\n", err.Error())
			return nil, fmt.Errorf("RemoveImage error: %s", err.Error())
		}
	}

	return &dockerpb.RemoveImageResponse{Message: "Removed: " + req.Id}, nil
}

// Networks
func (d *DockerService) GetNetworkList(ctx context.Context, req *dockerpb.GetNetworkListRequest) (*dockerpb.GetNetworkListResponse, error) {
	networkList, err := networks.GetNetworkList(ctx, d.cli)
	if err != nil {
		log.Printf("GetNetworkList error: %s\n", err.Error())
		return nil, fmt.Errorf("GetNetworkList error: %s", err.Error())
	}

	var serializedList dockerpb.GetNetworkListResponse
	serializedList.Items = make([]*dockerpb.NetworkBaseInfo, len(networkList))
	for i := range networkList {
		used, err := networks.IsNetworkUsed(ctx, d.cli, networkList[i].Id)
		if err != nil {
			log.Printf("GetNetworkList error: %s\n", err.Error())
			return nil, fmt.Errorf("GetNetworkList error: %s", err.Error())
		}
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
			IsUsed:        used,
		}
	}

	return &serializedList, nil
}

func (d *DockerService) GetNetworkByID(ctx context.Context, req *dockerpb.GetNetworkByIDRequest) (*dockerpb.GetNetworkByIDResponse, error) {
	network, err := networks.GetNetwork(ctx, d.cli, req.Id)
	if err != nil {
		fmt.Printf("GetNetwork error: %s\n", err.Error())
		return nil, err
	}

	used, err := networks.IsNetworkUsed(ctx, d.cli, req.Id)
	if err != nil {
		log.Printf("GetImageByID error: %s\n", err.Error())
		return nil, fmt.Errorf("GetImageByID error: %s", err.Error())
	}

	serializedNetwork := dockerpb.GetNetworkByIDResponse{
		Item: &dockerpb.NetworkBaseInfo{
			Name:          network.Name,
			Id:            network.Id,
			Driver:        network.Driver,
			EnableIpv6:    network.EnableIPv6,
			IpamDriver:    network.IPAMDriver,
			Subnet:        network.Subnet,
			Gateway:       network.Gateway,
			Attachable:    network.Attachable,
			DockerNetwork: network.DockerNetwork,
			IsUsed:        used,
		},
	}

	return &serializedNetwork, nil
}

func (d *DockerService) RemoveNetwork(ctx context.Context, req *dockerpb.RemoveNetworkRequest) (*dockerpb.RemoveNetworkResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("RemoveNetwork error: nil request")
	}

	used, err := networks.IsNetworkUsed(ctx, d.cli, req.Id)
	if err != nil {
		log.Printf("RemoveNetwork error: %s\n", err.Error())
		return nil, fmt.Errorf("RemoveNetwork error: %s", err.Error())
	}

	if used {
		return nil, fmt.Errorf("network is in use")
	}

	if err := networks.RemoveNetworkByID(ctx, d.cli, req.Id); err != nil {
		log.Printf("RemoveNetwork error: %s\n", err.Error())
		return nil, fmt.Errorf("RemoveNetwork error: %s", err.Error())
	}

	return &dockerpb.RemoveNetworkResponse{Message: "Removed: " + req.Id}, nil
}

// Volumes
func (d *DockerService) GetVolumeList(ctx context.Context, req *dockerpb.GetVolumeListRequest) (*dockerpb.GetVolumeListResponse, error) {
	volumeList, err := volumes.GetVolumeList(d.cli)
	if err != nil {
		log.Printf("GetVolumeList error: %s\n", err.Error())
		return nil, fmt.Errorf("GetVolumeList error: %s", err.Error())
	}

	var serializedList dockerpb.GetVolumeListResponse
	serializedList.Items = make([]*dockerpb.VolumeBaseInfo, len(volumeList))
	for i := range volumeList {
		used, err := containers.IsVolumeUsed(ctx, d.cli, volumeList[i].Name)
		if err != nil {
			log.Printf("GetVolumeList error: %s\n", err.Error())
			return nil, fmt.Errorf("GetVolumeList error: %s", err.Error())
		}
		serializedList.Items[i] = &dockerpb.VolumeBaseInfo{
			Name:       volumeList[i].Name,
			Driver:     volumeList[i].Driver,
			Mountpoint: volumeList[i].Mountpoint,
			Created:    volumeList[i].Created,
			IsUsed:     used,
		}
	}

	return &serializedList, nil
}

func (d *DockerService) GetVolumeByID(ctx context.Context, req *dockerpb.GetVolumeByIDRequest) (*dockerpb.GetVolumeByIDResponse, error) {
	volume, err := volumes.GetVolume(ctx, d.cli, req.Id)
	if err != nil {
		fmt.Printf("GetVolume error: %s\n", err.Error())
		return nil, err
	}

	used, err := containers.IsVolumeUsed(ctx, d.cli, volume.Name)
	if err != nil {
		log.Printf("GetVolumeByID error: %s\n", err.Error())
		return nil, fmt.Errorf("GetVolumeByID error: %s", err.Error())
	}

	serializedVolume := dockerpb.GetVolumeByIDResponse{
		Item: &dockerpb.VolumeBaseInfo{
			Name:       volume.Name,
			Driver:     volume.Driver,
			Mountpoint: volume.Mountpoint,
			Created:    volume.Created,
			IsUsed:     used,
		},
	}

	return &serializedVolume, nil
}

func (d *DockerService) RemoveVolume(ctx context.Context, req *dockerpb.RemoveVolumeRequest) (*dockerpb.RemoveVolumeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("RemoveVolume error: nil request")
	}

	if req.Force {
		if err := volumes.RemoveVolumeByID(ctx, d.cli, req.Id, req.Force); err != nil {
			log.Printf("RemoveVolume error: %s\n", err.Error())
			return nil, fmt.Errorf("RemoveVolume error: %s", err.Error())
		}
	} else {
		used, err := containers.IsVolumeUsed(ctx, d.cli, req.Id)
		if err != nil {
			log.Printf("RemoveVolume error: %s\n", err.Error())
			return nil, fmt.Errorf("RemoveVolume error: %s", err.Error())
		}

		if used {
			return nil, fmt.Errorf("volume is in use")
		}

		if err := volumes.RemoveVolumeByID(ctx, d.cli, req.Id, false); err != nil {
			log.Printf("RemoveVolume error: %s\n", err.Error())
			return nil, fmt.Errorf("RemoveVolume error: %s", err.Error())
		}
	}

	return &dockerpb.RemoveVolumeResponse{Message: "Removed: " + req.Id}, nil
}

// Deploy
func (d *DockerService) CreateFromGit(ctx context.Context, req *dockerpb.CreateFromGitRequest) (*dockerpb.CreateFromGitResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("CreateFromGit error: nil request")
	}

	newUUID := uuid.New().String()
	logCh := d.startBuildLogChannel(newUUID, 200)

	go d.startGitBuild(req, newUUID, logCh)

	return &dockerpb.CreateFromGitResponse{BuildId: newUUID}, nil
}

func (d *DockerService) StreamBuildLogs(req *dockerpb.StreamBuildLogsRequest, stream grpc.ServerStreamingServer[dockerpb.StreamBuildLogsResponse]) error {
	defer log.Println("StreamBuildLogs closed")
	ctx := stream.Context()

	d.logMu.RLock()
	ch, ok := d.logChans[req.BuildId]
	d.logMu.RUnlock()

	if !ok {
		return fmt.Errorf("build %s not found", req.BuildId)
	}

	for {
		select {
		case <-ctx.Done():
			{
				return nil
			}
		case line, ok := <-ch:
			{
				if !ok { //closed chan
					return nil
				}

				resp := &dockerpb.StreamBuildLogsResponse{
					Line: line,
				}
				if err := stream.Send(resp); err != nil {
					log.Printf("Send error: %s", err.Error())
					return err
				}
			}
		}
	}
}

// events
func (d *DockerService) GetContainerState(req *dockerpb.GetContainerStateRequest, stream grpc.ServerStreamingServer[dockerpb.GetContainerStateResponse]) error {
	ctx := stream.Context()

	eventCh := make(chan events.Message)
	errCh := make(chan error)
	go containers.EventStream(ctx, d.cli, eventCh, errCh)

	// log.Println("Start gRPC Server-Side streaming")

	for {
		select {
		case <-ctx.Done():
			{
				// log.Println("Stop gRPC Server-Side streaming")
				return nil
			}
		case err := <-errCh:
			{
				if err != nil {
					// log.Println("Stop gRPC Server-Side streaming")
					return err
				}
			}
		case event := <-eventCh:
			{
				// fmt.Println("New event in GetContainerState")
				resp := &dockerpb.GetContainerStateResponse{
					Type:   string(event.Type),
					Action: string(event.Action),
					Actor: &dockerpb.Actor{
						Id:         event.Actor.ID,
						Attributes: event.Actor.Attributes,
					},
				}
				if err := stream.Send(resp); err != nil {
					log.Printf("Send error: %s", err.Error())
					// log.Println("Stop gRPC Server-Side streaming")
					return err
				}
			}
		}
	}
}

func (d *DockerService) GetImageState(req *dockerpb.GetImageStateRequest, stream grpc.ServerStreamingServer[dockerpb.GetImageStateResponse]) error {
	ctx := stream.Context()

	eventCh := make(chan events.Message)
	errCh := make(chan error)
	go images.EventStream(ctx, d.cli, eventCh, errCh)

	// log.Println("Start gRPC Server-Side streaming")

	for {
		select {
		case <-ctx.Done():
			{
				// log.Println("Stop gRPC Server-Side streaming")
				return nil
			}
		case err := <-errCh:
			{
				if err != nil {
					// log.Println("Stop gRPC Server-Side streaming")
					return err
				}
			}
		case event := <-eventCh:
			{
				// fmt.Println("New event in GetContainerState")
				resp := &dockerpb.GetImageStateResponse{
					Type:   string(event.Type),
					Action: string(event.Action),
					Actor: &dockerpb.Actor{
						Id:         event.Actor.ID,
						Attributes: event.Actor.Attributes,
					},
				}
				if err := stream.Send(resp); err != nil {
					log.Printf("Send error: %s", err.Error())
					// log.Println("Stop gRPC Server-Side streaming")
					return err
				}
			}
		}
	}
}

func (d *DockerService) GetNetworkState(req *dockerpb.GetNetworkStateRequest, stream grpc.ServerStreamingServer[dockerpb.GetNetworkStateResponse]) error {
	ctx := stream.Context()

	eventCh := make(chan events.Message)
	errCh := make(chan error)
	go networks.EventStream(ctx, d.cli, eventCh, errCh)

	// log.Println("Start gRPC Server-Side streaming")

	for {
		select {
		case <-ctx.Done():
			{
				// log.Println("Stop gRPC Server-Side streaming")
				return nil
			}
		case err := <-errCh:
			{
				if err != nil {
					// log.Println("Stop gRPC Server-Side streaming")
					return err
				}
			}
		case event := <-eventCh:
			{
				// fmt.Println("New event in GetContainerState")
				resp := &dockerpb.GetNetworkStateResponse{
					Type:   string(event.Type),
					Action: string(event.Action),
					Actor: &dockerpb.Actor{
						Id:         event.Actor.ID,
						Attributes: event.Actor.Attributes,
					},
				}
				if err := stream.Send(resp); err != nil {
					log.Printf("Send error: %s", err.Error())
					// log.Println("Stop gRPC Server-Side streaming")
					return err
				}
			}
		}
	}
}

func (d *DockerService) GetVolumeState(req *dockerpb.GetVolumeStateRequest, stream grpc.ServerStreamingServer[dockerpb.GetVolumeStateResponse]) error {
	ctx := stream.Context()

	eventCh := make(chan events.Message)
	errCh := make(chan error)
	go volumes.EventStream(ctx, d.cli, eventCh, errCh)

	// log.Println("Start gRPC Server-Side streaming")

	for {
		select {
		case <-ctx.Done():
			{
				// log.Println("Stop gRPC Server-Side streaming")
				return nil
			}
		case err := <-errCh:
			{
				if err != nil {
					// log.Println("Stop gRPC Server-Side streaming")
					return err
				}
			}
		case event := <-eventCh:
			{
				// fmt.Println("New event in GetContainerState")
				resp := &dockerpb.GetVolumeStateResponse{
					Type:   string(event.Type),
					Action: string(event.Action),
					Actor: &dockerpb.Actor{
						Id:         event.Actor.ID,
						Attributes: event.Actor.Attributes,
					},
				}
				if err := stream.Send(resp); err != nil {
					log.Printf("Send error: %s", err.Error())
					// log.Println("Stop gRPC Server-Side streaming")
					return err
				}
			}
		}
	}
}

// common
func (d *DockerService) startBuildLogChannel(buildID string, bufSize int) chan string {
	ch := make(chan string, bufSize)

	d.logMu.Lock()
	d.logChans[buildID] = ch
	d.logMu.Unlock()

	return ch
}

func (d *DockerService) startGitBuild(req *dockerpb.CreateFromGitRequest, buildID string, logCh chan string) {
	ctx := context.Context(context.Background())
	defer func() {
		d.logMu.Lock()
		close(logCh)
		delete(d.logChans, buildID)
		d.logMu.Unlock()
	}()
	//
	//
	// Clonning repo from git into temp dir
	//
	//
	tempDir, err := os.MkdirTemp("", "docker-repo")
	if err != nil {
		common.SendLog(logCh, err.Error())
		log.Println(err.Error())
		return
	}

	createdDirLog := "Created tempDir: " + tempDir
	deletedDirLog := "Deleted tempDir: " + tempDir
	common.SendLog(logCh, createdDirLog)
	log.Println(createdDirLog)
	defer os.RemoveAll(tempDir)
	defer log.Println(deletedDirLog)

	if err := git.CloneRepo(&git.RepoCloner{}, req.Url, tempDir); err != nil {
		common.SendLog(logCh, err.Error())
		log.Println(err.Error())
		return
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
		if !common.FileExists(dockerfileDir) {
			notExist := "Dockerfile does not exist"
			common.SendLog(logCh, notExist)
			log.Println(notExist)
			return
		}
	} else {
		// If Dockerfile specified in the request, rewrite/create new Dockerfile
		if err := common.CreateFile(dockerfileDir, []byte(req.Dockerfile)); err != nil {
			common.SendLog(logCh, err.Error())
			log.Println(err.Error())
			return
		} else {
			dockerfileRedefined := "Dockerfile redefined"
			common.SendLog(logCh, dockerfileRedefined)
			log.Println(dockerfileRedefined)
		}
	}
	//
	//
	// Parse tags from repo (for image name)
	//
	//
	repo, err := common.RepoFromURL(req.Url)
	if err != nil {
		common.SendLog(logCh, err.Error())
		log.Println(err.Error())
		return
	}
	tags := []string{repo + ":latest"} //args.URL
	//
	//
	// Delete old containers
	//
	//
	if err := containers.RemoveContainerByImage(ctx, d.cli, tags[0]); err != nil {
		common.SendLog(logCh, err.Error())
		log.Println(err.Error())
		return
	}
	//
	//
	// Delete old images
	//
	//
	if err := images.RemoveImageByTag(ctx, d.cli, tags[0]); err != nil {
		common.SendLog(logCh, err.Error())
		log.Println(err.Error())
		return
	}
	//
	//
	// Build image
	//
	//
	newImageId, err := images.BuildImageNew(ctx, d.cli, tempDir, tags, logCh)
	if err != nil {
		common.SendLog(logCh, err.Error())
		log.Println(err.Error())
		return
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
		common.SendLog(logCh, err.Error())
		log.Println(err.Error())
		return
	}
	_, err = d.cli.BuildCachePrune(ctx, build.CachePruneOptions{All: true})
	if err != nil {
		common.SendLog(logCh, err.Error())
		log.Println(err.Error())
		return
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
			if !common.FileExists(dockercomposeDir) {
				common.SendLog(logCh, "created but not launched: "+newImageId)
				log.Println("Created image, but not command to launch container")
				return
			} else {
				if err := containers.ExecDockerComposeUp(ctx, dockercomposeDir, logCh); err != nil {
					common.SendLog(logCh, "the image was created, but an error occurred when starting the container (dockerCompose): "+newImageId)
					log.Println("Created image, but container not launched")
					return
				}
			}
		} else {
			if err := containers.ExecDockerRun(ctx, req.DockerRun, logCh); err != nil {
				common.SendLog(logCh, "the image was created, but an error occurred when starting the container (dockerRun): "+newImageId)
				log.Println("Created image, but container not launched")
				return
			}
		}
	} else {
		// If Dockerfile specified in the request, rewrite/create new Dockerfile
		if err := common.CreateFile(dockercomposeDir, []byte(req.DockerCompose)); err != nil {
			common.SendLog(logCh, err.Error())
			log.Println(err.Error())
			return
		} else {
			log.Println("DockerCompose redefined")
			if err := containers.ExecDockerComposeUp(ctx, dockercomposeDir, logCh); err != nil {
				common.SendLog(logCh, "the image was created, but an error occurred when starting the container (dockerCompose): "+newImageId)
				log.Println("Created image, but container not launched")
				return
			}
		}
	}

	common.SendLog(logCh, "created and launched: "+newImageId)
	log.Println("Created image and launched container")
}
