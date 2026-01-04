package ipc

import (
	"context"
	"log"
	"sync"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"google.golang.org/grpc"

	"fmt"

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
	containerList, err := containers.GetContainerBaseInfoList(d.cli)
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

func (d *DockerService) PauseContainerByID(ctx context.Context, req *dockerpb.PauseContainerByIDRequest) (*dockerpb.PauseContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("PauseContainerByID error: nil request")
	}

	if err := containers.PauseContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("PauseContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("PauseContainerByID error: %s", err.Error())
	}

	return &dockerpb.PauseContainerByIDResponse{Message: "Paused: " + req.Id}, nil
}

func (d *DockerService) UnpauseContainerByID(ctx context.Context, req *dockerpb.UnpauseContainerByIDRequest) (*dockerpb.UnpauseContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("UnpauseContainerByID error: nil request")
	}

	if err := containers.UnpauseContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("UnpauseContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("UnpauseContainerByID error: %s", err.Error())
	}

	return &dockerpb.UnpauseContainerByIDResponse{Message: "Unpaused: " + req.Id}, nil
}

func (d *DockerService) KillContainerByID(ctx context.Context, req *dockerpb.KillContainerByIDRequest) (*dockerpb.KillContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("KillContainerByID error: nil request")
	}

	if err := containers.KillContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("KillContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("KillContainerByID error: %s", err.Error())
	}

	return &dockerpb.KillContainerByIDResponse{Message: "Killed: " + req.Id}, nil
}

func (d *DockerService) StartContainerByID(ctx context.Context, req *dockerpb.StartContainerByIDRequest) (*dockerpb.StartContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("StartContainerByID error: nil request")
	}

	if err := containers.StartContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("StartContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("StartContainerByID error: %s", err.Error())
	}

	return &dockerpb.StartContainerByIDResponse{Message: "Started: " + req.Id}, nil
}

func (d *DockerService) StopContainerByID(ctx context.Context, req *dockerpb.StopContainerByIDRequest) (*dockerpb.StopContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("StopContainerByID error: nil request")
	}

	if err := containers.StopContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("StopContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("StopContainerByID error: %s", err.Error())
	}

	return &dockerpb.StopContainerByIDResponse{Message: "Stopped: " + req.Id}, nil
}

func (d *DockerService) RestartContainerByID(ctx context.Context, req *dockerpb.RestartContainerByIDRequest) (*dockerpb.RestartContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("RestartContainerByID error: nil request")
	}

	if err := containers.RestartContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("RestartContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("RestartContainerByID error: %s", err.Error())
	}

	return &dockerpb.RestartContainerByIDResponse{Message: "Restarted: " + req.Id}, nil
}

func (d *DockerService) RemoveContainerByID(ctx context.Context, req *dockerpb.RemoveContainerByIDRequest) (*dockerpb.RemoveContainerByIDResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("RemoveContainerByID error: nil request")
	}

	if err := containers.RemoveContainerByID(d.cli, req.Id); err != nil {
		fmt.Printf("RemoveContainerByID error: %s\n", err.Error())
		return nil, fmt.Errorf("RemoveContainerByID error: %s", err.Error())
	}

	return &dockerpb.RemoveContainerByIDResponse{Message: "Removed: " + req.Id}, nil
}

// Images
func (d *DockerService) GetImageList(ctx context.Context, req *dockerpb.GetImageListRequest) (*dockerpb.GetImageListResponse, error) {
	imageList, err := images.GetImageList(d.cli)
	if err != nil {
		fmt.Printf("GetImageList error: %s\n", err.Error())
		return nil, fmt.Errorf("GetImageList error: %s", err.Error())
	}

	var serializedList dockerpb.GetImageListResponse
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
func (d *DockerService) GetNetworkList(ctx context.Context, req *dockerpb.GetNetworkListRequest) (*dockerpb.GetNetworkListResponse, error) {
	networkList, err := networks.GetNetworkList(d.cli)
	if err != nil {
		fmt.Printf("GetNetworkList error: %s\n", err.Error())
		return nil, fmt.Errorf("GetNetworkList error: %s", err.Error())
	}

	var serializedList dockerpb.GetNetworkListResponse
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
func (d *DockerService) GetVolumeList(ctx context.Context, req *dockerpb.GetVolumeListRequest) (*dockerpb.GetVolumeListResponse, error) {
	volumeList, err := volumes.GetVolumeList(d.cli)
	if err != nil {
		fmt.Printf("GetVolumeList error: %s\n", err.Error())
		return nil, fmt.Errorf("GetVolumeList error: %s", err.Error())
	}

	var serializedList dockerpb.GetVolumeListResponse
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
					fmt.Println("Stop gRPC Server-Side streaming")
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
		fmt.Println(err.Error())
		return
	}

	createdDirLog := "Created tempDir: " + tempDir
	deletedDirLog := "Deleted tempDir: " + tempDir
	common.SendLog(logCh, createdDirLog)
	fmt.Println(createdDirLog)
	defer os.RemoveAll(tempDir)
	defer fmt.Println(deletedDirLog)

	if err := git.CloneRepo(&git.RepoCloner{}, req.Url, tempDir); err != nil {
		common.SendLog(logCh, err.Error())
		fmt.Println(err.Error())
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
			common.SendLog(logCh, err.Error())
			fmt.Println(err.Error())
			return
		}
	} else {
		// If Dockerfile specified in the request, rewrite/create new Dockerfile
		if err := common.CreateFile(dockerfileDir, []byte(req.Dockerfile)); err != nil {
			common.SendLog(logCh, err.Error())
			fmt.Println(err.Error())
			return
		} else {
			dockerfileRedefined := "Dockerfile redefined"
			common.SendLog(logCh, dockerfileRedefined)
			fmt.Println(dockerfileRedefined)
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
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
		return
	}
	//
	//
	// Delete old images
	//
	//
	if err := images.RemoveImageByTag(ctx, d.cli, tags[0]); err != nil {
		common.SendLog(logCh, err.Error())
		fmt.Println(err.Error())
		return
	}
	//
	//
	// Build image
	//
	//
	newImageId, err := images.BuildImageNew(d.cli, tempDir, tags, logCh)
	if err != nil {
		common.SendLog(logCh, err.Error())
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
		return
	}
	_, err = d.cli.BuildCachePrune(ctx, build.CachePruneOptions{All: true})
	if err != nil {
		common.SendLog(logCh, err.Error())
		fmt.Println(err.Error())
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
				fmt.Println("Created image, but not command to launch container")
				return
			} else {
				if err := containers.ExecDockerComposeUp(dockercomposeDir, logCh); err != nil {
					common.SendLog(logCh, "the image was created, but an error occurred when starting the container (dockerCompose): "+newImageId)
					fmt.Println("Created image, but container not launched")
					return
				}
			}
		} else {
			if err := containers.ExecDockerRun(req.DockerRun, logCh); err != nil {
				common.SendLog(logCh, "the image was created, but an error occurred when starting the container (dockerRun): "+newImageId)
				fmt.Println("Created image, but container not launched")
				return
			}
		}
	} else {
		// If Dockerfile specified in the request, rewrite/create new Dockerfile
		if err := common.CreateFile(dockercomposeDir, []byte(req.DockerCompose)); err != nil {
			common.SendLog(logCh, err.Error())
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println("DockerCompose redefined")
			if err := containers.ExecDockerComposeUp(dockercomposeDir, logCh); err != nil {
				common.SendLog(logCh, "the image was created, but an error occurred when starting the container (dockerCompose): "+newImageId)
				fmt.Println("Created image, but container not launched")
				return
			}
		}
	}

	common.SendLog(logCh, "created and launched: "+newImageId)
	fmt.Println("Created image and launched container")
}
